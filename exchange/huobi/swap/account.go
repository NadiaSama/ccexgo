package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type (
	TransferReq struct {
		data map[string]interface{}
	}

	TransferResp struct {
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Message string      `json:"string"`
		Success bool        `json:"success"`
	}

	transferAmount struct {
		amount float64
	}
)

const (
	TransferEndPoint    = "/v2/account/transfer"
	TransferSpotAccount = "spot"
	TransferSwapAccount = "swap"
)

//NewTransferReq build swap transfer reqeust amount will truncate 8 decimal places
func NewTransferReq(from, to, currency string, amount float64) *TransferReq {
	ret := TransferReq{
		data: make(map[string]interface{}),
	}

	ta := transferAmount{
		amount: amount,
	}

	ret.data["from"] = from
	ret.data["to"] = to
	ret.data["currency"] = currency
	ret.data["amount"] = &ta
	return &ret
}

func (tr *TransferReq) Serialize() ([]byte, error) {
	return json.Marshal(tr.data)
}

func (rc *RestClient) Transfer(ctx context.Context, req *TransferReq) (*TransferResp, error) {
	raw, err := req.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "serialize fail")
	}
	buf := bytes.NewBuffer(raw)
	var resp TransferResp
	if err := rc.RequestWithRawResp(ctx, http.MethodPost, TransferEndPoint, nil, buf, true, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		req, _ := req.Serialize()
		return nil, errors.Errorf("transfer fail resp='%+v', req='%s'", resp, string(req))
	}

	return &resp, nil
}

func (tf *transferAmount) MarshalJSON() ([]byte, error) {
	ret := strconv.FormatFloat(float64(tf.amount), 'f', 8, 64)
	return []byte(ret), nil
}
