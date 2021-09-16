package okex5

import (
	"context"
	"net/http"
	"net/url"
)

type (
	FundingRate struct {
		InstType        InstType `json:"InstType"`
		InstID          string   `json:"instId"`
		FundingRate     string   `json:"fundingRate"`
		NextFundingRate string   `json:"nextFundingRate"`
		FundingTime     string   `json:"fundingTime"`
		NextFundingTime string   `json:"nextFundingTime"`
	}
)

const (
	FundingEndPoint = "/api/v5/public/funding-rate"
)

func (rc *RestClient) FundingRate(ctx context.Context, instID string) ([]FundingRate, error) {
	var rates []FundingRate
	values := url.Values{}
	values.Add("instId", instID)
	if err := rc.Request(ctx, http.MethodGet, FundingEndPoint, values, nil, false, &rates); err != nil {
		return nil, err
	}

	return rates, nil
}
