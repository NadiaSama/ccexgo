package spot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
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

func (rc *RestClient) FeeRate(ctx context.Context, syms ...exchange.Symbol) ([]exchange.TradeFee, error) {
	param := url.Values{}
	if len(syms) != 0 {
		symbols := make([]string, len(syms))
		for i, sym := range syms {
			symbols[i] = sym.String()
		}
		param.Add("symbols", strings.Join(symbols, ","))
	}

	var result TradeFeeResult
	if err := rc.Request(ctx, http.MethodGet, "/v2/reference/transact-fee-rate", param,
		nil, true, &result); err != nil {
		return nil, err
	}
	if result.Code != huobi.CodeOK {
		b, _ := json.Marshal(result)
		return nil, huobi.NewError(string(b))
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

	var (
		mf    float64
		mfstr string
		tf    float64
		tfstr string
	)

	if len(tfr.ActualMakerFeeRate) != 0 {
		mfstr = tfr.ActualMakerFeeRate
	} else {
		mfstr = tfr.MakerFeeRate
	}
	mf, err = strconv.ParseFloat(mfstr, 64)
	if err != nil {
		return errors.WithMessagef(err, "parse float64 %s fail", mfstr)
	}

	if len(tfr.ActualTakerFeeRate) != 0 {
		tfstr = tfr.ActualTakerFeeRate
	} else {
		tfstr = tfr.TakerFeeRate
	}
	tf, err = strconv.ParseFloat(tfstr, 64)
	if err != nil {
		return errors.WithMessagef(err, "parse float64 %s fail", tfstr)
	}
	dst.Symbol = sym
	dst.Maker = mf
	dst.Taker = tf
	return nil
}
