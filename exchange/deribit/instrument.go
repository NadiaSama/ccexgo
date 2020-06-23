package deribit

import (
	"context"
	"strings"
)

type (
	InstrumentResult struct {
		TickSize            float64 `json:"tick_size"`
		TakerCommision      float64 `json:"taker_commision"`
		MakerCommision      float64 `json:"maker_commision"`
		Strike              float64 `json:"strike"`
		SettlementPeriod    string  `json:"settlement_period"`
		QuoteCurrency       string  `json:"quote_currency"`
		BaseCurreny         string  `json:"base_currency"`
		MinTradeAmount      float64 `json:"min_trade_amount"`
		Kind                string  `json:"kind"`
		IsActive            bool    `json:"is_active"`
		InstrumentName      string  `json:"instrument_name"`
		ExpirationTimestamp int64   `json:"expeiration_timestamp"`
		CreationTimestamp   int64   `json:"creation_timestamp"`
		ContractSize        float64 `json:"contract_size"`
	}
)

func (c *Client) OptionFetchInstruments(ctx context.Context, currency string) ([]InstrumentResult, error) {
	var ir []InstrumentResult
	param := map[string]interface{}{
		"currency": strings.ToUpper(currency),
		"kind":     "option",
		"expired":  false,
	}
	if err := c.call(ctx, "public/get_instruments", param, &ir, false); err != nil {
		return nil, err
	}

	return ir, nil
}
