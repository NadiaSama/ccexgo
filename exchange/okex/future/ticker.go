package future

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/shopspring/decimal"
)

type (
	tickerChannel struct {
		symbol exchange.FuturesSymbol
	}

	rawTicker struct {
		Last           decimal.Decimal `json:"last"`
		Open24H        decimal.Decimal `json:"open_24h"`
		BestBid        decimal.Decimal `json:"best_bid"`
		High24H        decimal.Decimal `json:"high_24h"`
		Low24H         decimal.Decimal `json:"low_24h"`
		Volume24H      decimal.Decimal `json:"volume_24h"`
		VolumeToken24H decimal.Decimal `json:"volume_token_24h"`
		BestAsk        decimal.Decimal `json:"best_ask"`
		OpenInterest   decimal.Decimal `json:"open_interest"`
		InstrumentID   string          `json:"instrument_id"`
		Timestamp      string          `json:"timestamp"`
		BestBidSize    decimal.Decimal `json:"best_bid_size"`
		BestAskSize    decimal.Decimal `json:"best_ask_size"`
		LastQty        decimal.Decimal `json:"last_qty"`
	}

	Ticker struct {
		Symbol         exchange.FuturesSymbol
		Last           decimal.Decimal
		LastQty        decimal.Decimal
		BestBid        decimal.Decimal
		BestBidSize    decimal.Decimal
		BestAsk        decimal.Decimal
		BestAskSize    decimal.Decimal
		Open24H        decimal.Decimal
		High24H        decimal.Decimal
		Low24H         decimal.Decimal
		Volume24H      decimal.Decimal
		VolumeToken24H decimal.Decimal
		OpenInterest   decimal.Decimal
		Time           time.Time
	}
)

const (
	TickerTable = "futures/ticker"
)

func init() {
	okex.SubscribeCB(TickerTable, parseTickerCB)
}

func NewTickerChannel(symbol exchange.FuturesSymbol) exchange.Channel {
	return &tickerChannel{
		symbol: symbol,
	}
}

func (tc *tickerChannel) String() string {
	return fmt.Sprintf("%s:%s", TickerTable, tc.symbol.String())
}

func parseTickerCB(table string, action string, raw json.RawMessage) (*rpc.Notify, error) {
	var rt []rawTicker
	if err := json.Unmarshal(raw, &rt); err != nil {
		return nil, err
	}

	r := rt[0]
	sym, err := ParseSymbol(r.InstrumentID)
	if err != nil {
		return nil, err
	}

	ts, err := okex.ParseTime(r.Timestamp)
	if err != nil {
		return nil, err
	}

	ticker := Ticker{
		Symbol:         sym,
		Last:           r.Last,
		LastQty:        r.LastQty,
		BestBid:        r.BestBid,
		BestBidSize:    r.BestBidSize,
		BestAsk:        r.BestAsk,
		BestAskSize:    r.BestAskSize,
		Open24H:        r.Open24H,
		High24H:        r.High24H,
		Low24H:         r.Low24H,
		Volume24H:      r.Volume24H,
		VolumeToken24H: r.VolumeToken24H,
		OpenInterest:   r.OpenInterest,
		Time:           ts,
	}

	return &rpc.Notify{
		Method: TickerTable,
		Params: &exchange.Ticker{
			LastPrice:   ticker.Last,
			BestBid:     ticker.BestBid,
			BestBidSize: ticker.BestBidSize,
			BestAsk:     ticker.BestAsk,
			BestAskSize: ticker.BestAskSize,
			Symbol:      ticker.Symbol,
			Time:        ticker.Time,
			Raw:         &ticker,
		},
	}, nil
}
