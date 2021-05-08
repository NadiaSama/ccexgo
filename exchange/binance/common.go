package binance

import (
	"fmt"
	"net/url"
)

func TradeParam(symbol string, st int64, et int64, fid int64, limit int) url.Values {
	value := url.Values{}
	value.Add("symbol", symbol)
	if st != 0 {
		value.Add("startTime", fmt.Sprintf("%d", st))
	}
	if et != 0 {
		value.Add("endTime", fmt.Sprintf("%d", et))
	}
	if fid != 0 {
		value.Add("fromId", fmt.Sprintf("%d", fid))
	}
	if limit != 0 {
		value.Add("limit", fmt.Sprintf("%d", limit))
	}
	return value
}
