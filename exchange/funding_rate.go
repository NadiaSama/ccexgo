package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	FundingRate struct {
		Symbol          Symbol
		FundingRate     decimal.Decimal
		NextFundingTime time.Time
		Time            time.Time
		Raw             interface{}
	}
)
