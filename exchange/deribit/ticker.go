package deribit

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	TickerGreeks struct {
		Vega  decimal.Decimal `json:"vega"`
		Theta decimal.Decimal `json:"theta"`
		Rho   decimal.Decimal `json:"rho"`
		Gamma decimal.Decimal `json:"gamma"`
		Delta decimal.Decimal `json:"delta"`
	}

	TickerStats struct {
		Volume      decimal.Decimal
		PriceChange decimal.Decimal `json:"price_change"`
		Low         decimal.Decimal
		High        decimal.Decimal
	}

	TickerResult struct {
		UnderlyingPrice        decimal.Decimal `json:"underlying_price"`
		UnderlyingIndex        string          `json:"underlying_index"`
		Timestamp              int64           `json:"timestamp"`
		State                  string          `json:"state"`
		Stats                  TickerStats     `json:"stats"`
		SettlementPrice        decimal.Decimal `json:"settlement_price"`
		OpenInterest           decimal.Decimal `json:"open_interest"`
		MinPrice               decimal.Decimal `json:"min_price"`
		MaxPrice               decimal.Decimal `json:"max_price"`
		MarkPrice              decimal.Decimal `json:"mark_price"`
		MarkIV                 decimal.Decimal `json:"mark_iv"`
		LastPrice              decimal.Decimal `json:"last_price"`
		InterestRate           decimal.Decimal `json:"interest_rate"`
		InstrumentName         string          `json:"instrument_name"`
		IndexPrice             decimal.Decimal `json:"index_price"`
		Greeks                 TickerGreeks    `json:"greeks"`
		EstimatedDeliveryPrice decimal.Decimal `json:"estimated_delivery_price"`
		BidIV                  decimal.Decimal `json:"bid_iv"`
		BestBidPrice           decimal.Decimal `json:"best_bid_price"`
		BestBidAmount          decimal.Decimal `json:"best_bid_amount"`
		BestAskPrice           decimal.Decimal `json:"best_ask_price"`
		BestAskAmount          decimal.Decimal `json:"best_ask_amount"`
		AskIV                  decimal.Decimal `json:"ask_iv"`
	}

	ChTicker struct {
		instrument string
	}
)

const (
	PublicTickerMethod = "public/ticker"
)

func init() {
	reigisterCB("ticker", parseNotifyTicker)
}

func NewTickerChannel(instrument string) *ChTicker {
	return &ChTicker{
		instrument: instrument,
	}
}

func (ct *ChTicker) String() string {
	return fmt.Sprintf("ticker.%s.100ms", ct.instrument)
}

func parseNotifyTicker(resp *Notify) (*rpc.Notify, error) {
	var tr TickerResult
	if err := json.Unmarshal(resp.Data, &tr); err != nil {
		return nil, errors.WithMessage(err, "unmarshal ticker result")
	}

	ticker, err := tr.Parse()
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: subscriptionMethod,
		Params: ticker,
	}, nil
}

func (tr *TickerResult) Parse() (*exchange.Ticker, error) {
	sym, err := ParseSymbol(tr.InstrumentName)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse instrument_name '%s'", tr.InstrumentName)
	}

	return &exchange.Ticker{
		Symbol:      sym,
		BestBid:     tr.BestBidPrice,
		BestBidSize: tr.BestBidAmount,
		BestAsk:     tr.BestAskPrice,
		BestAskSize: tr.BestAskAmount,
		Time:        tconv.Milli2Time(tr.Timestamp),
		Raw:         tr,
	}, nil
}
