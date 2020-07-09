package huobi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	TransactFeeRate struct {
		Symbol             string `json:"symbol"`
		MakerFeeRate       string `json:"makerFeeRate"`
		TakerFeeRate       string `json:"takerFeeRate"`
		ActualMakerFeeRate string `json:"actualMakerFeeRate"`
		ActualTakerFeeRate string `json:"actualTakerFeeRate"`
	}

	TradeFeeResult struct {
		Code    int               `json:"code"`
		Success bool              `json:"success"`
		Data    []TransactFeeRate `json:"data"`
	}
)

func (rc *RestClient) FeeRate(ctx context.Context, syms ...exchange.SpotSymbol) ([]exchange.TradeFee, error) {
	var param map[string]string
	if len(syms) != 0 {
		symbols := make([]string, len(syms))
		for i, sym := range syms {
			symbols[i] = sym.String()
		}
		param = map[string]string{}
		param["symbols"] = strings.Join(symbols, ",")
	}

	var result TradeFeeResult
	if err := rc.request(ctx, http.MethodGet, "/v2/reference/transact-fee-rate", param,
		nil, true, &result); err != nil {
		return nil, err
	}
	if result.Code != codeOK {
		b, _ := json.Marshal(result)
		return nil, newError(string(b))
	}

	ret := make([]exchange.TradeFee, len(result.Data))
	for i, tf := range result.Data {
		err := rc.parse(&tf, &(ret[i]))
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (rc *RestClient) parse(tfr *TransactFeeRate, dst *exchange.TradeFee) error {
	sym, err := rc.ParseSpotSymbol(tfr.Symbol)
	if err != nil {
		return err
	}

	mf, err := strconv.ParseFloat(tfr.MakerFeeRate, 64)
	if err != nil {
		return errors.WithMessagef(err, "parse float64 %s fail", tfr.MakerFeeRate)
	}
	tf, err := strconv.ParseFloat(tfr.TakerFeeRate, 64)
	if err != nil {
		return errors.WithMessagef(err, "parse float64 %s fail", tfr.TakerFeeRate)
	}

	dst.Symbol = sym
	dst.Maker = mf
	dst.Taker = tf
	return nil
}
