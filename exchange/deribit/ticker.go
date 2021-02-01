package deribit

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	TickerGreeks struct {
		Vega  decimal.Decimal
		Theta decimal.Decimal
		Rho   decimal.Decimal
		Gamma decimal.Decimal
		Delta decimal.Decimal
	}

	TickerStats struct {
		Volume      decimal.Decimal
		PriceChange decimal.Decimal
		Low         decimal.Decimal
		High        decimal.Decimal
	}

	TickerResult struct {
		UnderlyingPrice        decimal.Decimal
		UnderlyingIndex        string
		Timestamp              int64
		State                  string
		Stats                  TickerStats
		SettlementPrice        decimal.Decimal
		OpenInterest           decimal.Decimal
		MinPrice               decimal.Decimal
		MaxPrice               decimal.Decimal
		MarkPrice              decimal.Decimal
		MarkIV                 decimal.Decimal
		LastPrice              decimal.Decimal
		InterestRate           decimal.Decimal
		InstrumentName         string
		IndexPrice             decimal.Decimal
		Greeks                 TickerGreeks
		EstimatedDeliveryPrice decimal.Decimal
		BidIV                  decimal.Decimal
		BestBidPrice           decimal.Decimal
		BestBidAmount          decimal.Decimal
		BestAskPrice           decimal.Decimal
		BestAskAmount          decimal.Decimal
		AskIV                  decimal.Decimal
	}
)

const (
	PublicTickerMethod = "public/ticker"
)

func init() {
	reigisterCB("ticker", parseNotifyIndex)
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
	sym, err := ParseOptionSymbol(tr.InstrumentName)
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
		Raw:         &tr,
	}, nil
}
