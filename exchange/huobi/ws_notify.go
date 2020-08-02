package huobi

import (
	"encoding/json"

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

func (r *Response) Parse() (rpc.Response, error) {
	if r.Ping != 0 {
		return &rpc.Notify{
			Method: MethodPing,
			Params: r.Ping,
		}, nil
	}

	trades, err := ParseTrades(r.Tick)
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: r.Ch,
		Params: trades,
	}, nil
}
