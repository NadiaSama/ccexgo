package spot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	PlaceReq struct {
		data map[string]string
	}

	PlaceResp struct {
		Status  string `json:"status"`
		Data    string `json:"data"`
		ErrCode string `json:"err-code"`
		ErrMsg  string `json:"err-msg"`
	}
)

const (
	PlaceOrderEndPoint = "/v1/order/orders/place"
)

func NewPlaceReq(accountID string, symbol string, typ string, amount string) *PlaceReq {
	ret := &PlaceReq{
		data: make(map[string]string),
	}
	ret.data["account-id"] = accountID
	ret.data["symbol"] = symbol
	ret.data["type"] = typ
	ret.data["amount"] = amount

	return ret
}

func (pr *PlaceReq) Price(price string) *PlaceReq {
	pr.data["price"] = price
	return pr
}

func (pr *PlaceReq) Source(source string) *PlaceReq {
	pr.data["source"] = source
	return pr
}

func (pr *PlaceReq) ClientOrderID(coi string) *PlaceReq {
	pr.data["client-order-id"] = coi
	return pr
}

func (pr *PlaceReq) StopPrice(sp string) *PlaceReq {
	pr.data["stop-price"] = sp
	return pr
}

func (pr *PlaceReq) Operator(op string) *PlaceReq {
	pr.data["operator"] = op
	return pr
}

func (pr *PlaceReq) Serialize() ([]byte, error) {
	b, e := json.Marshal(pr.data)
	if e != nil {
		return nil, errors.WithMessage(e, "marshal json fail")
	}
	return b, nil
}

func (rc *RestClient) Place(ctx context.Context, req *PlaceReq) (*PlaceResp, error) {
	msg, err := req.Serialize()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(msg)
	var resp PlaceResp
	if err := rc.RequestWithRawResp(ctx, http.MethodPost, PlaceOrderEndPoint, nil, buf, true, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, errors.Errorf("place order error %+v", resp)
	}

	return &resp, nil
}

func (rc *RestClient) SubmitCancel(ctx context.Context, orderID string) (*PlaceResp, error) {
	url := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderID)

	var resp PlaceResp
	if err := rc.RequestWithRawResp(ctx, http.MethodPost, url, nil, nil, true, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, errors.Errorf("cancel order fail %+v", resp)
	}

	return &resp, nil
}
