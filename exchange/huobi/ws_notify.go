package huobi

import (
	"encoding/json"
	"errors"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Response struct {
		Ping    int             `json:"ping,omitempty"`
		Ch      string          `json:"ch,omitempty"`
		TS      int64           `json:"ts,omitempty"`
		Tick    json.RawMessage `json:"tick,omitempty"`
		ID      string          `json:"id,omitempty"`
		Status  string          `json:"status,omitempty"`
		Subbed  string          `json:"subbed,omitempty"`
		ErrCode string          `json:"err-code,omitempty"`
		ErrMsg  string          `json:"err-msg,omitempty"`
	}
)

var (
	//SkipError means the response can not be handled directly by Parse method
	SkipError = errors.New("skip error")
)

func (r *Response) Parse(raw []byte) (rpc.Response, error) {
	if r.Ping != 0 {
		return &rpc.Notify{
			Method: MethodPing,
			Params: r.Ping,
		}, nil
	}

	if r.ID != "" {
		return &rpc.Result{
			ID:     r.ID,
			Result: raw,
		}, nil
	}

	return nil, SkipError
}
