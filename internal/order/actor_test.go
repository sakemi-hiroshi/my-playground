package order

import (
	"sync"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- テスト用ヘルパー ---

type recorder struct {
	mu    sync.Mutex
	calls []string
}

func (r *recorder) record(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, s)
}

func (r *recorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.calls))
	copy(out, r.calls)
	return out
}

// --- スタブアクター ---

type stubCouponActor struct {
	applyResult interface{}
	rec         *recorder
}

func (s *stubCouponActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ApplyCoupon:
		_ = msg
		ctx.Respond(s.applyResult)
	case ReleaseCoupon:
		s.rec.record("coupon.release")
		ctx.Respond(CouponReleased{OrderID: msg.OrderID})
	}
}

type stubPointActor struct {
	useResult interface{}
	rec       *recorder
}

func (s *stubPointActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case UsePoint:
		_ = msg
		ctx.Respond(s.useResult)
	case RefundPoint:
		s.rec.record("point.refund")
		ctx.Respond(PointRefunded{OrderID: msg.OrderID})
	}
}

type stubPaymentActor struct {
	chargeResult interface{}
	rec          *recorder
}

func (s *stubPaymentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case Charge:
		_ = msg
		ctx.Respond(s.chargeResult)
	case Refund:
		s.rec.record("payment.refund")
		ctx.Respond(PaymentRefunded{OrderID: msg.OrderID})
	}
}

// stubCouponActorWithTimeouts は最初の N 回は応答せず（タイムアウト模擬）、その後 result を返す
type stubCouponActorWithTimeouts struct {
	timeoutCount int // 応答しない回数
	result       interface{}
	rec          *recorder
	callCount    int
}

func (s *stubCouponActorWithTimeouts) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ApplyCoupon:
		_ = msg
		if s.callCount < s.timeoutCount {
			s.callCount++
			// 応答しない（タイムアウトを模擬）
			return
		}
		ctx.Respond(s.result)
	case ReleaseCoupon:
		s.rec.record("coupon.release")
		ctx.Respond(CouponReleased{OrderID: msg.OrderID})
	}
}

// stubPointActorWithTimeouts は最初の N 回は応答せず（タイムアウト模擬）、その後 result を返す
type stubPointActorWithTimeouts struct {
	timeoutCount int
	result       interface{}
	rec          *recorder
	callCount    int
}

func (s *stubPointActorWithTimeouts) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case UsePoint:
		_ = msg
		if s.callCount < s.timeoutCount {
			s.callCount++
			return
		}
		ctx.Respond(s.result)
	case RefundPoint:
		s.rec.record("point.refund")
		ctx.Respond(PointRefunded{OrderID: msg.OrderID})
	}
}

// stubPaymentActorWithTimeouts は最初の N 回は応答せず（タイムアウト模擬）、その後 result を返す
type stubPaymentActorWithTimeouts struct {
	timeoutCount int
	result       interface{}
	rec          *recorder
	callCount    int
}

func (s *stubPaymentActorWithTimeouts) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case Charge:
		_ = msg
		if s.callCount < s.timeoutCount {
			s.callCount++
			return
		}
		ctx.Respond(s.result)
	case Refund:
		s.rec.record("payment.refund")
		ctx.Respond(PaymentRefunded{OrderID: msg.OrderID})
	}
}

// --- テスト ---

func TestOrderActor_SagaFlow(t *testing.T) {
	tests := []struct {
		name              string
		couponResult      interface{}
		pointResult       interface{}
		paymentResult     interface{}
		wantStatus        OrderStatus
		wantCompensations []string
	}{
		{
			name:              "全成功_completed",
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			pointResult:       PointUsed{OrderID: "o1", Amount: 100},
			paymentResult:     PaymentCompleted{OrderID: "o1", AmountYen: 900},
			wantStatus:        StatusCompleted,
			wantCompensations: []string{},
		},
		{
			name:              "クーポン失敗_補償なし_failed",
			couponResult:      CouponRejected{OrderID: "o1", Reason: "invalid coupon"},
			wantStatus:        StatusFailed,
			wantCompensations: []string{},
		},
		{
			name:              "ポイント失敗_クーポンのみ補償_failed",
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			pointResult:       PointRejected{OrderID: "o1", Reason: "insufficient points"},
			wantStatus:        StatusFailed,
			wantCompensations: []string{"coupon.release"},
		},
		{
			name:              "決済失敗_ポイントとクーポンを逆順補償_failed",
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			pointResult:       PointUsed{OrderID: "o1", Amount: 100},
			paymentResult:     PaymentFailed{OrderID: "o1", Reason: "card declined"},
			wantStatus:        StatusFailed,
			wantCompensations: []string{"point.refund", "coupon.release"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := actor.NewActorSystem()
			t.Cleanup(func() { system.Shutdown() })

			rec := &recorder{}
			couponPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &stubCouponActor{applyResult: tt.couponResult, rec: rec}
			}))
			pointPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &stubPointActor{useResult: tt.pointResult, rec: rec}
			}))
			paymentPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &stubPaymentActor{chargeResult: tt.paymentResult, rec: rec}
			}))

			orderPID := system.Root.Spawn(NewOrderActorProps(couponPID, pointPID, paymentPID))
			res, err := system.Root.RequestFuture(orderPID, StartOrder{
				OrderID:     "o1",
				CouponID:    "C100",
				PointAmount: 100,
				AmountYen:   1000,
			}, 5*time.Second).Result()

			require.NoError(t, err)
			got, ok := res.(OrderResult)
			require.True(t, ok)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.wantCompensations, rec.snapshot())
		})
	}
}

func TestOrderActor_RetryOnTimeout(t *testing.T) {
	// タイムアウトを短くしてテストを高速化
	orig := receiveTimeout
	receiveTimeout = 100 * time.Millisecond
	t.Cleanup(func() { receiveTimeout = orig })

	tests := []struct {
		name              string
		maxRetries        int
		couponTimeouts    int
		pointTimeouts     int
		paymentTimeouts   int
		couponResult      interface{}
		pointResult       interface{}
		paymentResult     interface{}
		wantStatus        OrderStatus
		wantCompensations []string
	}{
		{
			name:              "MaxRetries=2_1回タイムアウト後に成功_completed",
			maxRetries:        2,
			couponTimeouts:    1,
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			pointResult:       PointUsed{OrderID: "o1", Amount: 100},
			paymentResult:     PaymentCompleted{OrderID: "o1", AmountYen: 900},
			wantStatus:        StatusCompleted,
			wantCompensations: []string{},
		},
		{
			name:              "MaxRetries=2_リトライ上限超え_補償に移行_failed",
			maxRetries:        2,
			couponTimeouts:    3, // MaxRetries=2 を超える（0回目 + retry1 + retry2 で上限、3回目もタイムアウト）
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			wantStatus:        StatusFailed,
			wantCompensations: []string{},
		},
		{
			name:              "MaxRetries=0_タイムアウト即補償_failed",
			maxRetries:        0,
			couponTimeouts:    1,
			couponResult:      CouponApplied{OrderID: "o1", CouponID: "C100", DiscountYen: 100},
			wantStatus:        StatusFailed,
			wantCompensations: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := actor.NewActorSystem()
			t.Cleanup(func() { system.Shutdown() })

			rec := &recorder{}
			couponPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &stubCouponActorWithTimeouts{
					timeoutCount: tt.couponTimeouts,
					result:       tt.couponResult,
					rec:          rec,
				}
			}))
			pointPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				if tt.pointTimeouts > 0 {
					return &stubPointActorWithTimeouts{
						timeoutCount: tt.pointTimeouts,
						result:       tt.pointResult,
						rec:          rec,
					}
				}
				return &stubPointActor{useResult: tt.pointResult, rec: rec}
			}))
			paymentPID := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				if tt.paymentTimeouts > 0 {
					return &stubPaymentActorWithTimeouts{
						timeoutCount: tt.paymentTimeouts,
						result:       tt.paymentResult,
						rec:          rec,
					}
				}
				return &stubPaymentActor{chargeResult: tt.paymentResult, rec: rec}
			}))

			orderPID := system.Root.Spawn(NewOrderActorProps(couponPID, pointPID, paymentPID))
			res, err := system.Root.RequestFuture(orderPID, StartOrder{
				OrderID:     "o1",
				CouponID:    "C100",
				PointAmount: 100,
				AmountYen:   1000,
				MaxRetries:  tt.maxRetries,
			}, 5*time.Second).Result()

			require.NoError(t, err)
			got, ok := res.(OrderResult)
			require.True(t, ok)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.wantCompensations, rec.snapshot())
		})
	}
}
