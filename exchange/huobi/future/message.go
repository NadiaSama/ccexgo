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

	var resp huobi.Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errors.WithMessagef(err, "bad response '%s'", data)
	}

	fields := strings.Split(resp.Ch, ".")
	if len(fields) == 4 && fields[0] == "market" && fields[2] == "trade" && fields[3] == "detail" {
		code, ok := cc.codeMap[fields[1]]
		if !ok {
			return nil, errors.Errorf("bad response channel %s", resp.Ch)
		}
		f := []string{fields[0], code, fields[1], fields[2], fields[3]}
		resp.Ch = strings.Join(f, ".")
	}

	return resp.Parse()
}
