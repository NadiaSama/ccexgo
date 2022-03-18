package option

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	ListenKeyResp struct {
		ListenKey string `json:"listenKey"`
	}
)

const (
	ListenKeyEndPoint = "/vapi/v1/userDataStream"
)

func (rc *RestClient) GetListenKeyAddr(ctx context.Context) (string, error) {
	var ret ListenKeyResp
	resp := RestResp{
		Data: &ret,
	}
	if err := rc.Request(ctx, http.MethodPost, ListenKeyEndPoint, nil, nil, true, &resp); err != nil {
		return "", errors.WithMessage(err, "request listenKey fail")
	}

	return fmt.Sprintf("wss://%s/ws/%s", rc.wsAddr, ret.ListenKey), nil
}

func (rc *RestClient) PersistListenKey(ctx context.Context) error {
	var ret map[string]interface{}

	if err := rc.Request(ctx, http.MethodPut, ListenKeyEndPoint, nil, nil, true, &ret); err != nil {
		return errors.WithMessage(err, "persist listenKey fail")
	}
	return nil
}

func (rc *RestClient) DeleteListenKey(ctx context.Context) error {
	var ret map[string]interface{}

	if err := rc.Request(ctx, http.MethodDelete, ListenKeyEndPoint, nil, nil, true, &ret); err != nil {
		return errors.WithMessage(err, "delete listenKey fail")
	}
	return nil
}
