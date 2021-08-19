package swap

import (
	"context"

	"github.com/pkg/errors"
)

type (
	SwapFeeReq struct {
		ContractCode string `json:"contract_code"`
	}

	SwapFee struct {
		Symbol        string `json:"symbol"`
		ContractCode  string `json:"contract_code"`
		OpenMakerFee  string `json:"open_maker_fee"`
		OpenTakerFee  string `json:"open_taker_fee"`
		CloseMakerFee string `json:"close_maker_fee"`
		CloseTakerFee string `json:"close_taker_fee"`
		FeeAsset      string `json:"fee_asset"`
	}
)

const (
	SwapFeeEndPoint = "/swap-api/v1/swap_fee"
)

func NewSwapFeeReq(symbol string) *SwapFeeReq {
	return &SwapFeeReq{
		ContractCode: symbol,
	}
}

func (rc *RestClient) SwapFee(ctx context.Context, req *SwapFeeReq) ([]SwapFee, error) {
	var ret []SwapFee
	if err := rc.PrivatePostReq(ctx, SwapFeeEndPoint, req, &ret); err != nil {
		return nil, errors.WithMessage(err, "request swap_fee fail")
	}

	return ret, nil
}
