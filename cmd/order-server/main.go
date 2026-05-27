package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sakemi-hiroshi/my-playground/internal/order"
)

func main() {
	system := actor.NewActorSystem()
	order.RegisterSagaTracer(system)

	// 各サービスアクターは長命な1インスタンス（Singleton 的）
	couponPID := system.Root.Spawn(actor.PropsFromProducer(order.NewCouponActor))
	pointPID := system.Root.Spawn(actor.PropsFromProducer(order.NewPointActor))
	paymentPID := system.Root.Spawn(actor.PropsFromProducer(order.NewPaymentActor))

	h := order.NewOrderHandler(system, couponPID, pointPID, paymentPID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/orders", h.PostOrder)

	e.Logger.Fatal(e.Start(":8080"))
}
