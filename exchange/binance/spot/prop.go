package spot

import (
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

func (rc *RestClient) Property() exchange.Property {
	return exchange.Property{
		Trades: &exchange.TradesProp{
			MaxDuration: time.Hour * 168,
			SuportID:    true,
			SupportTime: true,
		},
	}
}
