package order

type (
	Charge struct {
		OrderID   string
		AmountYen int
		FailMode  FailMode
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
