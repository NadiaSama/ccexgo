package swap

import (
	"context"
	"encoding/json"

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

func NewPositionInfoRequest(code string) *PositionInfoRequest {
	return &PositionInfoRequest{
		ContractCode: code,
	}
}

func (pir *PositionInfoRequest) Serialize() ([]byte, error) {
	if pir.ContractCode != "" {
		r, err := json.Marshal(pir)
		return r, err
	}
	return nil, nil
}

func (rc *RestClient) PositionInfo(ctx context.Context, req *PositionInfoRequest) ([]Position, error) {
	var ret []Position
	if err := rc.PrivatePostReq(ctx, PositionInfoEndPoint, req, &ret); err != nil {
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
