package point

import "github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"

type (
	UsePoint struct {
		OrderID        string
		Amount         int
		FaultInjection faultinjection.FaultInjection
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
