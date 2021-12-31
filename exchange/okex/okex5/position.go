package okex5

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type (
	Positions struct {
		Adl         string `json:"adl"`
		AvailPos    string `json:"availPos"`
		AvgPx       string `json:"avgPx"`
		CTime       string `json:"cTime"`
		Ccy         string `json:"ccy"`
		DeltaBS     string `json:"deltaBS"`
		DeltaPA     string `json:"deltaPA"`
		GammaBS     string `json:"gammaBS"`
		GammaPA     string `json:"gammaPA"`
		IMR         string `json:"imr"`
		InstID      string `json:"instId"`
		InstType    string `json:"instType"`
		Interest    string `json:"interest"`
		Last        string `json:"last"`
		UsdPx       string `json:"usdPx"`
		Lever       string `json:"lever"`
		Liab        string `json:"liab"`
		LiabCcy     string `json:"liabCcy"`
		LiqPx       string `json:"liqPx"`
		MarkPx      string `json:"markPx"`
		Margin      string `json:"margin"`
		MgnMode     string `json:"mgnMode"`
		MgnRatio    string `json:"mgnRatio"`
		NMR         string `json:"nmr"`
		NotionalUsd string `json:"notionalUsd"`
		OptVal      string `json:"optVal"`
		PTime       string `json:"pTime"`
		Pos         string `json:"pos"`
		PosCcy      string `json:"posCcy"`
		PosID       string `json:"posId"`
		PosSide     string `json:"posSide"`
		ThetaBS     string `json:"thetaBS"`
		ThetaPA     string `json:"thetaPA"`
		TradeID     string `json:"tradeId"`
		UTime       string `json:"uTime"`
		Upl         string `json:"upl"`
		UplRatio    string `json:"uplRatio"`
		VegaBS      string `json:"vegaBS"`
		VegaPA      string `json:"vegaPA"`
	}

	PositionsReq struct {
		*GetRequest
	}
)

const (
	PositionsEndPoint = "/api/v5/account/positions"
)

func NewPositionsReq() *PositionsReq {
	return &PositionsReq{
		NewGetRequest(),
	}
}

func (pr *PositionsReq) InstType(typ string) *PositionsReq {
	pr.Add("instType", typ)
	return pr
}

func (pr *PositionsReq) InstID(id string) *PositionsReq {
	pr.Add("instId", id)
	return pr
}

func (pr *PositionsReq) PosID(id string) *PositionsReq {
	pr.Add("posId", id)
	return pr
}

func (rc *RestClient) Positions(ctx context.Context, req *PositionsReq) ([]Positions, error) {
	var ret []Positions
	if err := rc.Request(ctx, http.MethodGet, PositionsEndPoint, req.Values(), nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch positions fail")
	}

	return ret, nil
}
