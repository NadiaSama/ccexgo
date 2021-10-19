package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	ccexgo "github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Position struct {
		Symbol         string          `json:"symbol"`
		ContractCode   string          `json:"contract_code"`
		Volume         decimal.Decimal `json:"volume"`
		Available      decimal.Decimal `json:"available"`
		Frozen         decimal.Decimal `json:"frozen"`
		CostOpen       decimal.Decimal `json:"cost_open"`
		CostHold       decimal.Decimal `json:"cost_hold"`
		ProfitUnreal   decimal.Decimal `json:"profit_unreal"`
		ProfitRate     decimal.Decimal `json:"profit_rate"`
		LeverRate      int             `json:"lever_rate"`
		PositionMargin decimal.Decimal `json:"position_margin"`
		Direction      string          `json:"direction"`
		Profit         decimal.Decimal `json:"profit"`
		LastPrice      decimal.Decimal `json:"last_price"`
	}

	PositionInfoRequest struct {
		ContractCode string `json:"contract_code"`
	}
)

const (
	PositionInfoEndPoint = "/swap-api/v1/swap_position_info"
)

func (pir *PositionInfoRequest) Params() []byte {
	if pir.ContractCode != "" {
		r, _ := json.Marshal(pir)
		return r
	}
	return nil
}

func (rc *RestClient) PositionInfo(ctx context.Context, req *PositionInfoRequest) ([]Position, error) {
	var ret []Position
	raw := req.Params()
	var body io.Reader
	if len(raw) != 0 {
		body = bytes.NewBuffer(raw)
	}
	if err := rc.Request(ctx, http.MethodPost, PositionInfoEndPoint, nil, body, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get position info fail")
	}

	return ret, nil
}

func (p *Position) Transfer() (*ccexgo.Position, error) {
	sym, err := ParseSymbol(p.ContractCode)
	if err != nil {
		return nil, err
	}

	var side ccexgo.PositionSide
	if p.Direction == OrderDirectionBuy {
		side = ccexgo.PositionSideLong
	} else if p.Direction == OrderDirectionSell {
		side = ccexgo.PositionSideShort
	} else {
		return nil, errors.Errorf("unkown direction '%s'", p.Direction)
	}
	return &exchange.Position{
		Symbol:        sym,
		Side:          side,
		Mode:          ccexgo.PositionModeCross,
		Position:      p.Volume,
		AvailPosition: p.Available,
		AvgOpenPrice:  p.CostHold,
		UNRealizedPNL: p.Profit,
		Margin:        p.PositionMargin,
		Leverage:      decimal.NewFromInt(int64(p.LeverRate)),
		Raw:           p,
	}, nil
}
