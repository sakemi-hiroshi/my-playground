package order

import (
	"time"

	"github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"
)

type (
	FailModes struct {
		Coupon  faultinjection.FaultInjection
		Point   faultinjection.FaultInjection
		Payment faultinjection.FaultInjection
	}

	StartOrder struct {
		OrderID           string
		CouponID          string
		PointAmount       int
		AmountYen         int
		FailModes         FailModes
		RandomFailureRate float64
	}

	OrderStatus string

	OrderResult struct {
		OrderID string
		Status  OrderStatus
		Reason  string
	}

	// SagaEvent は EventStream に publish する観察用イベント
	SagaEvent struct {
		OrderID string
		Phase   string
		Detail  string
		At      time.Time
	}
)

const (
	StatusCompleted OrderStatus = "completed"
	StatusFailed    OrderStatus = "failed"
)
