package future

import (
	"encoding/json"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		*huobi.CodeC
		codeMap map[string]string
	}

	//Response format for huobi future websocket
	Response struct {
		Ping int             `json:"ping"`
		Ch   string          `json:"ch"`
		TS   int             `json:"ts"`
		Tick json.RawMessage `json:"tick"`
	}

	//callParam carry params which used by huobi websocket sub and pong
	callParam struct {
		Pong int    `json:"pong,omitempty"`
		Sub  string `json:"sub,omitempty"`
		ID   string `json:"id,omitempty"`
	}
)

func NewCodeC(cm map[string]string) *CodeC {
	c := make(map[string]string, len(cm))
	for k, v := range cm {
		c[k] = v
	}
	return &CodeC{
		codeMap: c,
		CodeC:   huobi.NewCodeC(),
	}
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	data, err := cc.Decompress(raw)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errors.WithMessagef(err, "bad response '%s'", data)
	}

	if resp.Ping != 0 {
		return &rpc.Notify{
			Method: huobi.MethodPing,
			Params: resp.Ping,
		}, nil
	}

	fields := strings.Split(resp.Ch, ".")
	if len(fields) != 4 || fields[0] != "market" || fields[2] != "trade" || fields[3] != "detail" {
		return nil, errors.Errorf("bad response channel %s", resp.Ch)
	}

	code, ok := cc.codeMap[fields[1]]
	if !ok {
		return nil, errors.Errorf("bad response channel %s", resp.Ch)
	}
	f := []string{fields[0], code, fields[1], fields[2], fields[3]}
	ch := strings.Join(f, ".")

	trades, err := huobi.ParseTrades(resp.Tick)
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: ch,
		Params: trades,
	}, nil
}
