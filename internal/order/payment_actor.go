package order

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type PaymentActor struct{}

func NewPaymentActor() actor.Actor {
	return &PaymentActor{}
}

func (a *PaymentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case Charge:
		a.handleCharge(ctx, msg)
	case Refund:
		ctx.Respond(PaymentRefunded{OrderID: msg.OrderID, AmountYen: msg.AmountYen})
	}
}

func (a *PaymentActor) handleCharge(ctx actor.Context, msg Charge) {
	switch msg.FailMode.Kind {
	case Fail:
		ctx.Respond(PaymentFailed{OrderID: msg.OrderID, Reason: "simulated failure"})
	case Delay:
		time.Sleep(msg.FailMode.Delay)
		ctx.Respond(PaymentCompleted{OrderID: msg.OrderID, AmountYen: msg.AmountYen})
	case Panic:
		panic("simulated panic")
	case NoReply:
		// 応答しない
	default:
		ctx.Respond(PaymentCompleted{OrderID: msg.OrderID, AmountYen: msg.AmountYen})
	}
}
