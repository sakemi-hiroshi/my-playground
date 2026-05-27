package order

import (
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
)

type FailModeRequest struct {
	Kind  string `json:"kind"`
	Delay string `json:"delay,omitempty"` // "100ms", "1s" など
}

type FailModesRequest struct {
	Coupon  FailModeRequest `json:"coupon"`
	Point   FailModeRequest `json:"point"`
	Payment FailModeRequest `json:"payment"`
}

type OrderRequest struct {
	OrderID           string           `json:"order_id"`
	CouponID          string           `json:"coupon_id"`
	PointAmount       int              `json:"point_amount"`
	AmountYen         int              `json:"amount_yen"`
	FailModes         FailModesRequest `json:"fail_modes"`
	RandomFailureRate float64          `json:"random_failure_rate"`
}

type OrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
}

type OrderHandler struct {
	system     *actor.ActorSystem
	couponPID  *actor.PID
	pointPID   *actor.PID
	paymentPID *actor.PID
}

func NewOrderHandler(
	system *actor.ActorSystem,
	couponPID, pointPID, paymentPID *actor.PID,
) *OrderHandler {
	return &OrderHandler{
		system:     system,
		couponPID:  couponPID,
		pointPID:   pointPID,
		paymentPID: paymentPID,
	}
}

func (h *OrderHandler) PostOrder(c echo.Context) error {
	var req OrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// OrderActor を受注ごとに spawn（per-request actor）
	pid := h.system.Root.Spawn(NewOrderActorProps(h.couponPID, h.pointPID, h.paymentPID))

	// HTTP は同期境界なので RequestFuture で完了を待つ
	fut := h.system.Root.RequestFuture(pid, StartOrder{
		OrderID:           req.OrderID,
		CouponID:          req.CouponID,
		PointAmount:       req.PointAmount,
		AmountYen:         req.AmountYen,
		FailModes:         toOrderFailModes(req.FailModes),
		RandomFailureRate: req.RandomFailureRate,
	}, 10*time.Second)

	res, err := fut.Result()
	if err != nil {
		return echo.NewHTTPError(http.StatusGatewayTimeout, "order processing timeout")
	}

	result := res.(OrderResult)
	return c.JSON(http.StatusOK, OrderResponse{
		OrderID: result.OrderID,
		Status:  string(result.Status),
		Reason:  result.Reason,
	})
}

func toOrderFailModes(r FailModesRequest) FailModes {
	return FailModes{
		Coupon:  toFailMode(r.Coupon),
		Point:   toFailMode(r.Point),
		Payment: toFailMode(r.Payment),
	}
}

func toFailMode(r FailModeRequest) FailMode {
	d, _ := time.ParseDuration(r.Delay)
	return FailMode{
		Kind:  Kind(r.Kind),
		Delay: d,
	}
}
