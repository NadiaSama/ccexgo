package exchange

import "time"

type (
	//TradesProp specific property which used to build Trades request
	TradesProp struct {
		MaxDuration time.Duration
		SuportID    bool
		SupportTime bool
	}

	FinanceProp struct {
		MaxDuration time.Duration
		SuportID    bool
		SupportTime bool
	}

	Property struct {
		Trades  *TradesProp
		Finance *FinanceProp
	}
)
