package order

type (
	UsePoint struct {
		OrderID  string
		Amount   int
		FailMode FailMode
	}

	PointUsed struct {
		OrderID string
		Amount  int
	}

	PointRejected struct {
		OrderID string
		Reason  string
	}

	RefundPoint struct {
		OrderID string
		Amount  int
	}

	PointRefunded struct {
		OrderID string
		Amount  int
	}
)
