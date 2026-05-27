package coupon

import "github.com/sakemi-hiroshi/my-playground/internal/failmode"

type (
	ApplyCoupon struct {
		OrderID  string
		CouponID string
		FailMode failmode.FailMode
	}

	CouponApplied struct {
		OrderID     string
		CouponID    string
		DiscountYen int
	}

	CouponRejected struct {
		OrderID string
		Reason  string
	}

	ReleaseCoupon struct {
		OrderID  string
		CouponID string
	}

	CouponReleased struct {
		OrderID string
	}
)
