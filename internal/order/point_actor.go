package order

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
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
	switch msg.FailMode.Kind {
	case Fail:
		ctx.Respond(PointRejected{OrderID: msg.OrderID, Reason: "simulated failure"})
	case Delay:
		time.Sleep(msg.FailMode.Delay)
		ctx.Respond(PointUsed{OrderID: msg.OrderID, Amount: msg.Amount})
	case Panic:
		panic("simulated panic")
	case NoReply:
		// 応答しない
	default:
		ctx.Respond(PointUsed{OrderID: msg.OrderID, Amount: msg.Amount})
	}
}
