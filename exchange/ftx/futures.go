package ftx

import (
	"context"
	"fmt"
	"net/http"
)

type (
	FutureInfo struct {
		Ask                 float64 `json:"ask"`
		Bid                 float64 `json:"bid"`
		Change1H            float64 `json:"change1h"`
		Change24H           float64 `json:"change24h"`
		ChangeBod           float64 `json:"changeBod"`
		VolumeUsd24h        float64 `json:"volumeUsd24h"`
		Volume              float64 `json:"volume"`
		Description         string  `json:"description"`
		Enabled             bool    `json:"enabled"`
		Expired             bool    `json:"expired"`
		Expiry              string  `json:"expiry"`
		Index               float64 `json:"index"`
		ImfFactor           float64 `json:"imfFactor"`
		Last                float64 `json:"last"`
		LowerBound          float64 `json:"lowerBound"`
		Mark                float64 `json:"mark"`
		Name                string  `json:"name"`
		Perpetual           bool    `json:"perpetual"`
		PositionLimitWtight float64 `json:"positionLimitWeight"`
		PostOnly            bool    `json:"postOnly"`
		PriceIncrement      float64 `json:"priceIncrement"`
		SizeIncrement       float64 `json:"sizeIncrement"`
		Underlying          string  `json:"underlying"`
		UpperBound          float64 `json:"upperBound"`
		Type                string  `json:"type"`
	}
)

func (rc *RestClient) Future(ctx context.Context, sym string) (*FutureInfo, error) {
	var info FutureInfo
	path := fmt.Sprintf("/futures/%s", sym)
	if err := rc.request(ctx, http.MethodGet, path, nil, nil, false, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (rc *RestClient) Futures(ctx context.Context) ([]FutureInfo, error) {
	var infos []FutureInfo
	if err := rc.request(ctx, http.MethodGet, "/futures", nil, nil, false, &infos); err != nil {
		return nil, err
	}
	return infos, nil
}
