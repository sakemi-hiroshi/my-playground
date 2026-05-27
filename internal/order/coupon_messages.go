package order

type (
	ApplyCoupon struct {
		OrderID        string
		CouponID       string
		FaultInjection FaultInjection
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
