package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	apphttp "github.com/sakemi-hiroshi/my-playground/internal/http"
	"github.com/sakemi-hiroshi/my-playground/internal/service/coupon"
	"github.com/sakemi-hiroshi/my-playground/internal/service/payment"
	"github.com/sakemi-hiroshi/my-playground/internal/service/point"
	"github.com/sakemi-hiroshi/my-playground/internal/tracing"
)

func main() {
	system := actor.NewActorSystem()
	tracing.RegisterSagaTracer(system)

	// 各サービスアクターは長命な1インスタンス（Singleton 的）
	couponPID := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))
	pointPID := system.Root.Spawn(actor.PropsFromProducer(point.NewPointActor))
	paymentPID := system.Root.Spawn(actor.PropsFromProducer(payment.NewPaymentActor))

	h := apphttp.NewOrderHandler(system, couponPID, pointPID, paymentPID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/orders", h.PostOrder)

	e.Logger.Fatal(e.Start(":8080"))
}
