package order

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
)

func RegisterSagaTracer(system *actor.ActorSystem) {
	system.EventStream.Subscribe(func(evt interface{}) {
		e, ok := evt.(SagaEvent)
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
