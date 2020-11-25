package swap

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	rawTicker struct {
		InstrumentID   string          `json:"instrument_id"`
		Last           decimal.Decimal `json:"last"`
		LastQty        decimal.Decimal `json:"last_qty"`
		BestAsk        decimal.Decimal `json:"best_ask"`
		BestAskSize    decimal.Decimal `json:"best_ask_size"`
		BestBid        decimal.Decimal `json:"best_bid"`
		BestBidSize    decimal.Decimal `json:"best_bid_size"`
		Open24H        decimal.Decimal `json:"open_24h"`
		High24H        decimal.Decimal `json:"high_24h"`
		Low24H         decimal.Decimal `json:"low_24h"`
		Volume24H      decimal.Decimal `json:"volume_24h"`
		VolumeToken24H decimal.Decimal `json:"volume_token_24h"`
		OpenInterest   decimal.Decimal `json:"open_interest"`
		Timestamp      string          `json:"timestamp"`
	}

	Ticker struct {
		Symbol         exchange.SwapSymbol
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

	TickerChannel struct {
		symbol exchange.SwapSymbol
	}
)

const (
	tickerTable = "swap/ticker"
)

func NewTickerChannel(sym exchange.SwapSymbol) exchange.Channel {
	return &TickerChannel{
		symbol: sym,
	}
}

func (tc *TickerChannel) String() string {
	return fmt.Sprintf("%s:%s", tickerTable, tc.symbol.String())
}

func init() {
	okex.SubscribeCB(tickerTable, parseTickerCB)
}
func parseTickerCB(table string, action string, raw json.RawMessage) (*rpc.Notify, error) {
	var rt []rawTicker
	if err := json.Unmarshal(raw, &rt); err != nil {
		return nil, err
	}

	r := rt[0]
	ts, err := okex.ParseTime(r.Timestamp)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse timestamp '%s'", r.Timestamp)
	}

	fields := strings.Split(r.InstrumentID, "-")
	if len(fields) != 3 {
		return nil, errors.Errorf("bad symbol '%s'", r.InstrumentID)
	}
	sym := okex.NewSwapSymbol(fmt.Sprintf("%s-%s", fields[0], fields[1]))

	return &rpc.Notify{
		Method: table,
		Params: &Ticker{
			Symbol:         sym,
			Time:           ts,
			Last:           r.Last,
			LastQty:        r.LastQty,
			BestAsk:        r.BestAsk,
			BestAskSize:    r.BestAskSize,
			BestBid:        r.BestBid,
			BestBidSize:    r.BestBidSize,
			Open24H:        r.Open24H,
			Volume24H:      r.Volume24H,
			VolumeToken24H: r.VolumeToken24H,
			OpenInterest:   r.OpenInterest,
		},
	}, nil
}
