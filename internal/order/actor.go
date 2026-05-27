package order

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

const receiveTimeout = 3 * time.Second

type OrderActor struct {
	behavior actor.Behavior

	// DI で受け取るサービスの PID
	couponPID  *actor.PID
	pointPID   *actor.PID
	paymentPID *actor.PID

	// Saga 実行中の状態
	msg     StartOrder
	replyTo *actor.PID
	comps   compensations
}

func NewOrderActorProps(couponPID, pointPID, paymentPID *actor.PID) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		a := &OrderActor{
			couponPID:  couponPID,
			pointPID:   pointPID,
			paymentPID: paymentPID,
		}
		a.behavior.Become(a.initial)
		return a
	})
}

func (a *OrderActor) Receive(ctx actor.Context) {
	a.behavior.Receive(ctx)
}

// --- 状態: initial ---

func (a *OrderActor) initial(ctx actor.Context) {
	msg, ok := ctx.Message().(StartOrder)
	if !ok {
		return
	}
	a.msg = msg
	a.replyTo = ctx.Sender()

	ctx.Request(a.couponPID, ApplyCoupon{
		OrderID:  msg.OrderID,
		CouponID: msg.CouponID,
		FailMode: msg.FailModes.Coupon,
	})
	ctx.SetReceiveTimeout(receiveTimeout)
	a.behavior.Become(a.awaitingCoupon)

	ctx.ActorSystem().EventStream.Publish(SagaEvent{
		OrderID: msg.OrderID, Phase: "started", At: time.Now(),
	})
}

// --- 状態: awaitingCoupon ---

func (a *OrderActor) awaitingCoupon(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case CouponApplied:
		ctx.CancelReceiveTimeout()
		a.comps.push(func(ctx actor.Context) {
			ctx.Request(a.couponPID, ReleaseCoupon{OrderID: a.msg.OrderID, CouponID: a.msg.CouponID})
		})
		ctx.Request(a.pointPID, UsePoint{
			OrderID:  a.msg.OrderID,
			Amount:   a.msg.PointAmount,
			FailMode: a.msg.FailModes.Point,
		})
		ctx.SetReceiveTimeout(receiveTimeout)
		a.behavior.Become(a.awaitingPoint)

		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "coupon_applied",
			Detail: msg.CouponID, At: time.Now(),
		})

	case CouponRejected:
		ctx.CancelReceiveTimeout()
		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "failed",
			Detail: msg.Reason, At: time.Now(),
		})
		a.finish(ctx, StatusFailed, msg.Reason)

	case *actor.ReceiveTimeout:
		a.startCompensation(ctx, "timeout waiting CouponApplied")
	}
}

// --- 状態: awaitingPoint ---

func (a *OrderActor) awaitingPoint(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case PointUsed:
		ctx.CancelReceiveTimeout()
		a.comps.push(func(ctx actor.Context) {
			ctx.Request(a.pointPID, RefundPoint{OrderID: a.msg.OrderID, Amount: a.msg.PointAmount})
		})
		ctx.Request(a.paymentPID, Charge{
			OrderID:   a.msg.OrderID,
			AmountYen: a.msg.AmountYen,
			FailMode:  a.msg.FailModes.Payment,
		})
		ctx.SetReceiveTimeout(receiveTimeout)
		a.behavior.Become(a.awaitingCharge)

		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "point_used",
			Detail: "ok", At: time.Now(),
		})

	case PointRejected:
		ctx.CancelReceiveTimeout()
		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "compensating",
			Detail: msg.Reason, At: time.Now(),
		})
		a.startCompensation(ctx, msg.Reason)

	case *actor.ReceiveTimeout:
		a.startCompensation(ctx, "timeout waiting PointUsed")
	}
}

// --- 状態: awaitingCharge ---

func (a *OrderActor) awaitingCharge(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case PaymentCompleted:
		ctx.CancelReceiveTimeout()
		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "completed",
			Detail: "ok", At: time.Now(),
		})
		a.finish(ctx, StatusCompleted, "")
		_ = msg

	case PaymentFailed:
		ctx.CancelReceiveTimeout()
		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "compensating",
			Detail: msg.Reason, At: time.Now(),
		})
		a.startCompensation(ctx, msg.Reason)

	case *actor.ReceiveTimeout:
		a.startCompensation(ctx, "timeout waiting PaymentCompleted")
	}
}

// --- 状態: compensating ---

func (a *OrderActor) compensating(ctx actor.Context) {
	switch ctx.Message().(type) {
	case CouponReleased, PointRefunded, PaymentRefunded:
		if a.comps.len() > 0 {
			a.comps.dispatchNext(ctx)
			return
		}
		ctx.ActorSystem().EventStream.Publish(SagaEvent{
			OrderID: a.msg.OrderID, Phase: "failed",
			Detail: "compensation done", At: time.Now(),
		})
		a.finish(ctx, StatusFailed, "compensation completed")
	}
}

// --- ヘルパー ---

func (a *OrderActor) startCompensation(ctx actor.Context, reason string) {
	if a.comps.len() == 0 {
		a.finish(ctx, StatusFailed, reason)
		return
	}
	a.comps.dispatchNext(ctx)
	a.behavior.Become(a.compensating)
}

func (a *OrderActor) finish(ctx actor.Context, status OrderStatus, reason string) {
	if a.replyTo != nil {
		ctx.Send(a.replyTo, OrderResult{
			OrderID: a.msg.OrderID,
			Status:  status,
			Reason:  reason,
		})
	}
	ctx.Stop(ctx.Self())
}
