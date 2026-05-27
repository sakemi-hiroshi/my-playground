package order

import (
	"time"
)

type (
	FailModes struct {
		Coupon  FailMode
		Point   FailMode
		Payment FailMode
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
