package point

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"
)

type PointActor struct{}

func NewPointActor() actor.Actor {
	return &PointActor{}
}

func (a *PointActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case UsePoint:
		a.handleUse(ctx, msg)
	case RefundPoint:
		ctx.Respond(PointRefunded{OrderID: msg.OrderID, Amount: msg.Amount})
	}
}

func (a *PointActor) handleUse(ctx actor.Context, msg UsePoint) {
	switch msg.FaultInjection.Kind {
	case faultinjection.Fail:
		ctx.Respond(PointRejected{OrderID: msg.OrderID, Reason: "simulated failure"})
	case faultinjection.Delay:
		time.Sleep(msg.FaultInjection.Delay)
		ctx.Respond(PointUsed{OrderID: msg.OrderID, Amount: msg.Amount})
	case faultinjection.Panic:
		panic("simulated panic")
	case faultinjection.NoReply:
		// 応答しない
	default:
		ctx.Respond(PointUsed{OrderID: msg.OrderID, Amount: msg.Amount})
	}
}
