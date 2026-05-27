package coupon

import "github.com/sakemi-hiroshi/my-playground/internal/order/faultinjection"

type (
	ApplyCoupon struct {
		OrderID        string
		CouponID       string
		FaultInjection faultinjection.FaultInjection
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
