package coupon_test

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/failmode"
	"github.com/sakemi-hiroshi/my-playground/internal/service/coupon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCouponActor_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name        string
		msg         coupon.ApplyCoupon
		wantType    interface{}
		wantTimeout bool
		timeout     time.Duration
	}{
		{
			name:     "正系_CouponAppliedが返る",
			msg:      coupon.ApplyCoupon{OrderID: "ord-1", CouponID: "C100"},
			wantType: coupon.CouponApplied{},
			timeout:  time.Second,
		},
		{
			name:     "FailMode=fail_CouponRejectedが返る",
			msg:      coupon.ApplyCoupon{OrderID: "ord-2", CouponID: "C100", FailMode: failmode.FailMode{Kind: failmode.Fail}},
			wantType: coupon.CouponRejected{},
			timeout:  time.Second,
		},
		{
			name:     "FailMode=delay_遅延後にCouponAppliedが返る",
			msg:      coupon.ApplyCoupon{OrderID: "ord-3", CouponID: "C100", FailMode: failmode.FailMode{Kind: failmode.Delay, Delay: 100 * time.Millisecond}},
			wantType: coupon.CouponApplied{},
			timeout:  time.Second,
		},
		{
			name:        "FailMode=noreply_タイムアウトになる",
			msg:         coupon.ApplyCoupon{OrderID: "ord-4", CouponID: "C100", FailMode: failmode.FailMode{Kind: failmode.NoReply}},
			wantTimeout: true,
			timeout:     200 * time.Millisecond,
		},
		{
			name:        "FailMode=panic_応答なしになる",
			msg:         coupon.ApplyCoupon{OrderID: "ord-5", CouponID: "C100", FailMode: failmode.FailMode{Kind: failmode.Panic}},
			wantTimeout: true,
			timeout:     500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := actor.NewActorSystem()
			t.Cleanup(func() { system.Shutdown() })

			pid := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))
			res, err := system.Root.RequestFuture(pid, tt.msg, tt.timeout).Result()

			if tt.wantTimeout {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.IsType(t, tt.wantType, res)
		})
	}
}

func TestCouponActor_ApplyCoupon_正系の詳細(t *testing.T) {
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	pid := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))
	res, err := system.Root.RequestFuture(pid, coupon.ApplyCoupon{OrderID: "ord-1", CouponID: "C100"}, time.Second).Result()

	require.NoError(t, err)
	got, ok := res.(coupon.CouponApplied)
	require.True(t, ok)
	assert.Equal(t, "ord-1", got.OrderID)
	assert.Equal(t, "C100", got.CouponID)
	assert.Positive(t, got.DiscountYen)
}

func TestCouponActor_panic後にRestartして動作すること(t *testing.T) {
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	pid := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))

	// panic させる（応答なしを期待）
	_, err := system.Root.RequestFuture(pid, coupon.ApplyCoupon{
		OrderID: "ord-5", CouponID: "C100", FailMode: failmode.FailMode{Kind: failmode.Panic},
	}, 500*time.Millisecond).Result()
	assert.Error(t, err)

	// Restart 後に正常動作すること
	res, err := system.Root.RequestFuture(pid, coupon.ApplyCoupon{OrderID: "ord-retry", CouponID: "C100"}, time.Second).Result()
	require.NoError(t, err)
	assert.IsType(t, coupon.CouponApplied{}, res)
}

func TestCouponActor_ReleaseCoupon(t *testing.T) {
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	pid := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))
	res, err := system.Root.RequestFuture(pid, coupon.ReleaseCoupon{OrderID: "ord-6", CouponID: "C100"}, time.Second).Result()

	require.NoError(t, err)
	got, ok := res.(coupon.CouponReleased)
	require.True(t, ok)
	assert.Equal(t, "ord-6", got.OrderID)
}
