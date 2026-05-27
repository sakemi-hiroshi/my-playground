package point

import "github.com/sakemi-hiroshi/my-playground/internal/failmode"

type (
	UsePoint struct {
		OrderID  string
		Amount   int
		FailMode failmode.FailMode
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
