package binance

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

func ParseTimestamp(ts int64) time.Time {
	return time.Unix(ts/1000, ts%1000*1e6)
}

func ToTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func ToTradeID(raw interface{}) (int64, error) {
	if raw == nil {
		return 0, nil
	}

	switch t := raw.(type) {
	case int:
	case int16:
	case int32:
	case int64:
		return int64(t), nil
	}
	return 0, errors.Errorf("invalid type %+v", raw)
}

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
