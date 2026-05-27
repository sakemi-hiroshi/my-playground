package payment

import "github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"

type (
	Charge struct {
		OrderID        string
		AmountYen      int
		FaultInjection faultinjection.FaultInjection
	}

	PaymentCompleted struct {
		OrderID   string
		AmountYen int
	}

	PaymentFailed struct {
		OrderID string
		Reason  string
	}

	Refund struct {
		OrderID   string
		AmountYen int
	}

	PaymentRefunded struct {
		OrderID   string
		AmountYen int
	}
)
