package payment

import "github.com/sakemi-hiroshi/my-playground/internal/failmode"

type (
	Charge struct {
		OrderID   string
		AmountYen int
		FailMode  failmode.FailMode
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
