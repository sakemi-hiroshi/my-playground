package payment

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"
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
	switch msg.FaultInjection.Kind {
	case faultinjection.Fail:
		ctx.Respond(PaymentFailed{OrderID: msg.OrderID, Reason: "simulated failure"})
	case faultinjection.Delay:
		time.Sleep(msg.FaultInjection.Delay)
		ctx.Respond(PaymentCompleted{OrderID: msg.OrderID, AmountYen: msg.AmountYen})
	case faultinjection.Panic:
		panic("simulated panic")
	case faultinjection.NoReply:
		// 応答しない
	default:
		ctx.Respond(PaymentCompleted{OrderID: msg.OrderID, AmountYen: msg.AmountYen})
	}
}
