package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
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

	//Response format for huobi
	Response struct {
		Ping int             `json:"ping"`
		Ch   string          `json:"ch"`
		TS   int             `json:"ts"`
		Tick json.RawMessage `json:"tick"`
	}
	responseParseCB func(*Response) (*rpc.Notify, error)
)

const (
	huobiPING = "ping"
)

func (CodeC) Encode(req rpc.Request) ([]byte, error) {
	return json.Marshal(&Request{
		Sub: req.Method(),
		ID:  fmt.Sprintf("id%d", req.ID().Num),
	})
}

//Decode huobi response current
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
			Method: huobiPING,
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
	fields[1] = code
	ch := strings.Join(fields, ".")

	trades, err := parseTrades(resp.Tick)
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: ch,
		Params: trades,
	}, nil
}
