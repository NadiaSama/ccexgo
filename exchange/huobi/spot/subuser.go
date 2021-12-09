package spot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	SubUserListReq struct {
		*exchange.RestReq
	}

	SubUserData struct {
		UID       int    `json:"uid"`
		UserState string `json:"userState"`
	}

	SubUserListResp struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		NextID  int           `json:"nextId"`
		Data    []SubUserData `json:"data"`
	}

	SubUserAccountReq struct {
		uid int
	}
)

const (
	SubUserListEndPoint = "/v2/sub-user/user-list"
)

func NewSubUserListReq() *SubUserListReq {
	return &SubUserListReq{
		RestReq: exchange.NewRestReq(),
	}
}

func (sr *SubUserListReq) FromID(id int) *SubUserListReq {
	sr.RestReq.AddFields("fromId", id)
	return sr
}

func (rc *RestClient) SubUserList(ctx context.Context, req *SubUserListReq) (*SubUserListResp, error) {
	var ret SubUserListResp

	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "build request param fail")
	}

	if err := rc.RequestWithRawResp(ctx, http.MethodGet, SubUserListEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	if ret.Message != "" {
		return nil, errors.Errorf("request fail code: %d msg: %s", ret.Code, ret.Message)
	}

	return &ret, nil
}

func NewSubUserAccountReq(uid int) *SubUserAccountReq {

	return &SubUserAccountReq{
		uid: uid,
	}
}

func (rc *RestClient) SubUserAccount(ctx context.Context, req *SubUserAccountReq) ([]BalanceResp, error) {
	var ret []BalanceResp

	uri := fmt.Sprintf("/v1/account/accounts/%d", req.uid)
	if err := rc.RestClient.Request(ctx, http.MethodGet, uri, nil, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request account fail")
	}

	return ret, nil
}
