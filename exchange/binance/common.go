package binance

import (
	"fmt"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

func ParseTimestamp(ts int64) time.Time {
	return time.Unix(ts/1000, ts%1000*1e6)
}
func TradeParam(req *exchange.TradeReqParam) url.Values {
	value := url.Values{}
	value.Add("symbol", req.Symbol.String())
	if !req.StartTime.IsZero() {
		value.Add("startTime", fmt.Sprintf("%d", req.StartTime.UnixNano()/1e6))
	}
	if !req.EndTime.IsZero() {
		value.Add("endTime", fmt.Sprintf("%d", req.EndTime.UnixNano()/1e6))
	}
	if req.StartID != "" {
		value.Add("fromId", req.StartID)
	}
	if req.Limit != 0 {
		value.Add("limit", fmt.Sprintf("%d", req.Limit))
	}
	return value
}
