package order

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
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
	switch msg.FailMode.Kind {
	case Fail:
		ctx.Respond(CouponRejected{OrderID: msg.OrderID, Reason: "simulated failure"})
	case Delay:
		time.Sleep(msg.FailMode.Delay)
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	case Panic:
		panic("simulated panic")
	case NoReply:
		// 応答しない
	default:
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	}
}
