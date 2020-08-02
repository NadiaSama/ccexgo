package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		//map huobi websocket subscribe futures symbol to contract info
		codeMap map[string]string
		decoder *gzip.Reader
	}

	Request struct {
		Sub string `json:"sub"`
		ID  string `json:"id"`
	}

	Pong struct {
		Pong int `json:"pong"`
	}

	//Response format for huobi websocket
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
	responseParseCB func(*Response) (*rpc.Notify, error)
)

const (
	methodPING     = "ping"
	methodPONG     = "pong"
	methodSubscibe = "sub"
)

func NewCodeC(cm map[string]string) *CodeC {
	c := make(map[string]string, len(cm))
	for k, v := range cm {
		c[k] = v
	}
	return &CodeC{
		codeMap: c,
	}
}

func (CodeC) Encode(req rpc.Request) ([]byte, error) {
	cm := req.Params().(*callParam)
	return json.Marshal(cm)
}

//Decode huobi response current only futures Trade and pong is support
func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	buf := bytes.NewBuffer(raw)
	if cc.decoder == nil {
		reader, err := gzip.NewReader(buf)
		if err != nil {
			return nil, errors.WithMessagef(err, "create decompress reader fail")
		}
		cc.decoder = reader
	} else {
		if err := cc.decoder.Reset(buf); err != nil {
			cc.decoder = nil
			return nil, errors.WithMessagef(err, "reset decompress reader fail")
		}
	}
	msg, err := ioutil.ReadAll(cc.decoder)
	if err != nil {
		return nil, errors.WithMessagef(err, "read decompress data fail")
	}

	var resp Response
	if err := json.Unmarshal(msg, &resp); err != nil {
		return nil, errors.WithMessagef(err, "bad response %s", string(raw))
	}

	if resp.Ping != 0 {
		return &rpc.Notify{
			Method: methodPING,
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

	trades, err := parseTrades(resp.Tick)
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: ch,
		Params: trades,
	}, nil
}
