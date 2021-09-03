package ftx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	Candle struct {
		Open      float64
		Close     float64
		High      float64
		Low       float64
		Volume    float64
		Time      float64
		StartTime string `json:"startTime"`
	}

	CandleReq struct {
		markName   string
		resolution int
		startTime  int64
		endTime    int64
		limit      int
	}
)

const (
	CandlesLimit = 1000
)

func NewCandelReq(name string, resolution int) *CandleReq {
	return &CandleReq{
		markName:   name,
		resolution: resolution,
	}
}

func (cr *CandleReq) StartTime(secs int64) *CandleReq {
	cr.startTime = secs
	return cr
}

func (cr *CandleReq) EndTime(secs int64) *CandleReq {
	cr.endTime = secs
	return cr
}

func (cr *CandleReq) Limit(l int) *CandleReq {
	cr.limit = l
	return cr
}

//Candles fetch ftx candles in ascending order
func (rc *RestClient) Candles(ctx context.Context, cr *CandleReq) ([]Candle, error) {
	var ret []Candle

	endPoint := fmt.Sprintf("/markets/%s/candles", cr.markName)
	values := url.Values{}
	values.Add("resolution", fmt.Sprintf("%d", cr.resolution))
	if cr.startTime != 0 {
		values.Add("start_time", fmt.Sprintf("%d", cr.startTime))
	}
	if cr.endTime != 0 {
		values.Add("end_time", fmt.Sprintf("%d", cr.endTime))
	}
	if cr.limit != 0 {
		values.Add("limit", fmt.Sprintf("%d", cr.limit))
	}

	if err := rc.request(ctx, http.MethodGet, endPoint, values, nil, false, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch candles fail")
	}

	return ret, nil
}

//Klines fetch klines in reverse orders
func (rc *RestClient) Klines(ctx context.Context, kr *exchange.KlineReq) ([]exchange.Kline, error) {
	if kr.Symbol == nil {
		return nil, errors.Errorf("missing symbol")
	}

	var ret []exchange.Kline
	secs := kr.Resolution.Secs()
	if secs == 0 {
		secs = int(kr.Resolution)
	}
	req := NewCandelReq(kr.Symbol.String(), secs)

	if !kr.StartTime.IsZero() {
		req.StartTime(kr.StartTime.Unix())
	}
	if !kr.EndTime.IsZero() {
		req.EndTime(kr.EndTime.Unix())
	}

	total := kr.Limit
	if total > CandlesLimit {
		req.Limit(CandlesLimit)
	} else if total > 0 {
		req.Limit(total)
	} else if total == 0 {
		req.Limit(CandlesLimit)
	}

	for {
		var (
			lastEndTime time.Time
		)
		candles, err := rc.Candles(ctx, req)
		if err != nil {
			return ret, errors.WithMessage(err, "request candles fail")
		}

		sort.Slice(candles, func(i, j int) bool {
			ci := candles[i]
			cj := candles[j]
			return ci.Time > cj.Time
		})

		for i := range candles {
			c := candles[i]
			k := c.Transform(kr.Symbol)
			ret = append(ret, *k)

			if total != 0 && len(ret) >= total {
				return ret, nil
			}
			lastEndTime = k.Time
		}

		req.EndTime(lastEndTime.Unix() - 1)
	}
}

func (c *Candle) Transform(symbol exchange.Symbol) *exchange.Kline {
	ts := int64(c.Time)
	t := time.Unix(ts/1000, ts%1000)
	return &exchange.Kline{
		Symbol: symbol,
		Open:   c.Open,
		Close:  c.Close,
		High:   c.High,
		Low:    c.Low,
		Volume: c.Volume,
		Time:   t,
		Raw:    c,
	}
}
