package order

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCouponActor_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name        string
		msg         ApplyCoupon
		wantType    interface{}
		wantTimeout bool
		timeout     time.Duration
	}{
		{
			name:     "正系_CouponAppliedが返る",
			msg:      ApplyCoupon{OrderID: "ord-1", CouponID: "C100"},
			wantType: CouponApplied{},
			timeout:  time.Second,
		},
		{
			name:     "FaultInjection=fail_CouponRejectedが返る",
			msg:      ApplyCoupon{OrderID: "ord-2", CouponID: "C100", FaultInjection: FaultInjection{Kind: Fail}},
			wantType: CouponRejected{},
			timeout:  time.Second,
		},
		{
			name:     "FaultInjection=delay_遅延後にCouponAppliedが返る",
			msg:      ApplyCoupon{OrderID: "ord-3", CouponID: "C100", FaultInjection: FaultInjection{Kind: Delay, Delay: 100 * time.Millisecond}},
			wantType: CouponApplied{},
			timeout:  time.Second,
		},
		{
			name:        "FaultInjection=noreply_タイムアウトになる",
			msg:         ApplyCoupon{OrderID: "ord-4", CouponID: "C100", FaultInjection: FaultInjection{Kind: NoReply}},
			wantTimeout: true,
			timeout:     200 * time.Millisecond,
		},
		{
			name:        "FaultInjection=panic_応答なしになる",
			msg:         ApplyCoupon{OrderID: "ord-5", CouponID: "C100", FaultInjection: FaultInjection{Kind: Panic}},
			wantTimeout: true,
			timeout:     500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := actor.NewActorSystem()
			t.Cleanup(func() { system.Shutdown() })

			pid := system.Root.Spawn(actor.PropsFromProducer(NewCouponActor))
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

	pid := system.Root.Spawn(actor.PropsFromProducer(NewCouponActor))
	res, err := system.Root.RequestFuture(pid, ApplyCoupon{OrderID: "ord-1", CouponID: "C100"}, time.Second).Result()

	require.NoError(t, err)
	got, ok := res.(CouponApplied)
	require.True(t, ok)
	assert.Equal(t, "ord-1", got.OrderID)
	assert.Equal(t, "C100", got.CouponID)
	assert.Positive(t, got.DiscountYen)
}

func TestCouponActor_panic後にRestartして動作すること(t *testing.T) {
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	pid := system.Root.Spawn(actor.PropsFromProducer(NewCouponActor))

	// panic させる（応答なしを期待）
	_, err := system.Root.RequestFuture(pid, ApplyCoupon{
		OrderID: "ord-5", CouponID: "C100", FaultInjection: FaultInjection{Kind: Panic},
	}, 500*time.Millisecond).Result()
	assert.Error(t, err)

	// Restart 後に正常動作すること
	res, err := system.Root.RequestFuture(pid, ApplyCoupon{OrderID: "ord-retry", CouponID: "C100"}, time.Second).Result()
	require.NoError(t, err)
	assert.IsType(t, CouponApplied{}, res)
}

func TestCouponActor_ReleaseCoupon(t *testing.T) {
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	pid := system.Root.Spawn(actor.PropsFromProducer(NewCouponActor))
	res, err := system.Root.RequestFuture(pid, ReleaseCoupon{OrderID: "ord-6", CouponID: "C100"}, time.Second).Result()

	require.NoError(t, err)
	got, ok := res.(CouponReleased)
	require.True(t, ok)
	assert.Equal(t, "ord-6", got.OrderID)
}
