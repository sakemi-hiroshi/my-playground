package coupon

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/failmode"
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
	case failmode.Fail:
		ctx.Respond(CouponRejected{OrderID: msg.OrderID, Reason: "simulated failure"})
	case failmode.Delay:
		time.Sleep(msg.FailMode.Delay)
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	case failmode.Panic:
		panic("simulated panic")
	case failmode.NoReply:
		// 応答しない
	default:
		ctx.Respond(CouponApplied{OrderID: msg.OrderID, CouponID: msg.CouponID, DiscountYen: 100})
	}
}
