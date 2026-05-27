package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	apphttp "github.com/sakemi-hiroshi/my-playground/internal/http"
	"github.com/sakemi-hiroshi/my-playground/internal/order"
	"github.com/sakemi-hiroshi/my-playground/internal/service/coupon"
	"github.com/sakemi-hiroshi/my-playground/internal/service/payment"
	"github.com/sakemi-hiroshi/my-playground/internal/service/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*actor.ActorSystem, *actor.PID, *actor.PID, *actor.PID) {
	t.Helper()
	system := actor.NewActorSystem()
	t.Cleanup(func() { system.Shutdown() })

	couponPID := system.Root.Spawn(actor.PropsFromProducer(coupon.NewCouponActor))
	pointPID := system.Root.Spawn(actor.PropsFromProducer(point.NewPointActor))
	paymentPID := system.Root.Spawn(actor.PropsFromProducer(payment.NewPaymentActor))
	return system, couponPID, pointPID, paymentPID
}

func TestOrderHandler_PostOrder(t *testing.T) {
	tests := []struct {
		name      string
		body      apphttp.OrderRequest
		wantOrder order.OrderStatus
	}{
		{
			name: "正系_completed",
			body: apphttp.OrderRequest{
				OrderID:     "o1",
				CouponID:    "C100",
				PointAmount: 100,
				AmountYen:   1000,
			},
			wantOrder: order.StatusCompleted,
		},
		{
			name: "決済失敗_failed",
			body: apphttp.OrderRequest{
				OrderID:   "o2",
				CouponID:  "C100",
				AmountYen: 1000,
				FailModes: apphttp.FailModesRequest{
					Payment: apphttp.FailModeRequest{Kind: "fail"},
				},
			},
			wantOrder: order.StatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system, couponPID, pointPID, paymentPID := setup(t)
			h := apphttp.NewOrderHandler(system, couponPID, pointPID, paymentPID)

			e := echo.New()
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			require.NoError(t, h.PostOrder(c))

			assert.Equal(t, http.StatusOK, rec.Code)
			var res apphttp.OrderResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&res))
			assert.Equal(t, string(tt.wantOrder), res.Status)
		})
	}
}
