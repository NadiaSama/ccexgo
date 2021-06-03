package okex5

import (
	"context"
	"net/http"
	"net/url"
)

type (
	Instrument struct {
		InstType  InstType `json:"instType"`
		InstID    string   `json:"instId"`
		Uly       string   `json:"uly"`
		Category  string   `json:"category"`
		BaseCcy   string   `json:"baseCcy"`
		QuoteCcy  string   `json:"quoteCcy"`
		SettleCcy string   `json:"settleCcy"`
		CtVal     string   `json:"CtVal"`
		CtMul     string   `json:"CtMul"`
		CtValCcy  string   `json:"ctValCcy"`
		OptType   string   `json:"optType"`
		Stk       string   `json:"stk"`
		ListTime  string   `json:"listTime"`
		ExpTime   string   `json:"expTime"`
		Lever     string   `json:"lever"`
		TickSz    string   `json:"tickSz"`
		LotSz     string   `json:"lotSz"`
		MinSz     string   `json:"minSz"`
		CtType    string   `json:"ctType"`
		Alias     string   `json:"alias"`
		State     string   `json:"state"`
	}
)

const (
	InstrumentEndPoint = "/api/v5/public/instruments"
)

func (rc *RestClient) Instruments(ctx context.Context, typ InstType) ([]Instrument, error) {
	var ret []Instrument
	v := url.Values{}
	v.Add("instType", string(typ))
	if err := rc.Request(ctx, http.MethodGet, InstrumentEndPoint, v, nil, false, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
