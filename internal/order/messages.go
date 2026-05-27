package order

import (
	"time"
)

type (
	FailModes struct {
		Coupon  FaultInjection
		Point   FaultInjection
		Payment FaultInjection
	}

	StartOrder struct {
		OrderID           string
		CouponID          string
		PointAmount       int
		AmountYen         int
		FailModes         FailModes
		RandomFailureRate float64
		MaxRetries        int // 0 = リトライなし（既存挙動と同じ）
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
