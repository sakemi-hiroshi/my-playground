package tracing

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/sakemi-hiroshi/my-playground/internal/order"
)

func RegisterSagaTracer(system *actor.ActorSystem) {
	system.EventStream.Subscribe(func(evt interface{}) {
		e, ok := evt.(order.SagaEvent)
		if !ok {
			return
		}
		slog.Info("[saga]",
			"order_id", e.OrderID,
			"phase", e.Phase,
			"detail", e.Detail,
		)
	})
}
