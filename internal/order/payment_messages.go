package order

type (
	Charge struct {
		OrderID        string
		AmountYen      int
		FaultInjection FaultInjection
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
