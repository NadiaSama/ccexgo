package huobi

import (
	"encoding/json"
	"errors"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Response struct {
		Ping int             `json:"ping"`
		Ch   string          `json:"ch"`
		TS   int             `json:"ts"`
		Tick json.RawMessage `json:"tick"`
	}
)

var (
	//SkipError means the response can not be handled directly by Parse method
	SkipError = errors.New("skip error")
)

func (r *Response) Parse() (rpc.Response, error) {
	if r.Ping != 0 {
		return &rpc.Notify{
			Method: MethodPing,
			Params: r.Ping,
		}, nil
	}

	return nil, SkipError
}
