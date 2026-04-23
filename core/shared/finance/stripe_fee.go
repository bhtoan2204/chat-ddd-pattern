package finance

import (
	"fmt"

	"wechat-clone/core/shared/pkg/stackErr"
)

type StripeFeePolicy struct {
	RateBPS    int64
	FlatAmount int64
}

func (p StripeFeePolicy) Compute(amount int64) (int64, error) {
	switch {
	case amount <= 0:
		return 0, stackErr.Error(fmt.Errorf("amount must be greater than 0"))
	case p.RateBPS < 0:
		return 0, stackErr.Error(fmt.Errorf("rate_bps must be greater than or equal to 0"))
	case p.FlatAmount < 0:
		return 0, stackErr.Error(fmt.Errorf("flat_amount must be greater than or equal to 0"))
	}

	rateFee := int64(0)
	if p.RateBPS > 0 {
		rateFee = (amount*p.RateBPS + 9999) / 10000
	}

	totalFee := rateFee + p.FlatAmount
	if totalFee < 0 {
		return 0, stackErr.Error(fmt.Errorf("stripe fee overflow for amount=%d", amount))
	}

	return totalFee, nil
}
