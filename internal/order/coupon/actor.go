package coupon

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"
)

type CouponActor struct{}

func NewCouponActor() actor.Actor {
	return &CouponActor{}
}

func (a *CouponActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ApplyCoupon:
		a.handleApply(ctx, msg)
	case ReleaseCoupon:
		ctx.Respond(CouponReleased{OrderID: msg.OrderID})
	}
}

func (a *CouponActor) handleApply(ctx actor.Context, msg ApplyCoupon) {
	switch msg.FaultInjection.Kind {
	case faultinjection.Fail:
		ctx.Respond(CouponRejected{OrderID: msg.OrderID, Reason: "simulated failure"})
	case faultinjection.Delay:
		time.Sleep(msg.FaultInjection.Delay)
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	case faultinjection.Panic:
		panic("simulated panic")
	case faultinjection.NoReply:
		// 応答しない
	default:
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	}
}
