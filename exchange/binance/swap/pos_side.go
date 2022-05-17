package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
)

type (
	GetPositionSideRequest struct {
		*binance.RestReq
	}

	GetPositionSideResp struct {
		binance.APIError
		DualSidePosition bool `json:"dualSidePosition"`
	}

	SetPositionSideRequest struct {
		*binance.RestReq
		dualSide bool
	}

	SetPositionSideResp struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
)

const (
	PositionSidePath = "/fapi/v1/positionSide/dual"
)

func (rc *RestClient) GetPositionSide(ctx context.Context, req *GetPositionSideRequest) (*GetPositionSideResp, error) {
	if rc.side != nil {
		return rc.side, nil
	}

	var side GetPositionSideResp
	if err := rc.GetRequest(ctx, PositionSidePath, req, true, &side); err != nil {
		return nil, errors.WithMessage(err, "get dual position side fail")
	}

	rc.side = &side
	return &side, nil
}

func (rc *RestClient) SetPositionSide(ctx context.Context, req *SetPositionSideRequest) error {
	values, err := req.Values()
	if err != nil {
		return errors.WithMessage(err, "get request values fail")
	}

	var resp SetPositionSideResp
	if err := rc.Request(ctx, http.MethodPost, PositionSidePath, values, nil, true, &resp); err != nil {
		return errors.WithMessage(err, "set position side fail")
	}

	if resp.Code != 200 {
		return errors.Errorf("set position fail resp=%+v", resp)
	}

	rc.side = &GetPositionSideResp{
		DualSidePosition: req.dualSide,
	}

	return nil
}

func NewGetPositionSideRequest() *GetPositionSideRequest {
	return &GetPositionSideRequest{
		binance.NewRestReq(),
	}
}

func NewSetPositionSideRequest(dualSide bool) *SetPositionSideRequest {
	req := binance.NewRestReq()
	req.AddFields("dualSidePosition", dualSide)
	return &SetPositionSideRequest{
		RestReq:  req,
		dualSide: dualSide,
	}
}
