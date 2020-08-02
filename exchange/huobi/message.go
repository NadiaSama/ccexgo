package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		decoder *gzip.Reader
	}
)

const (
	MethodPing     = "ping"
	MethodPong     = "pong"
	MethodSubscibe = "sub"
)

func NewCodeC() *CodeC {
	return &CodeC{
		decoder: nil,
	}
}

func (CodeC) Encode(req rpc.Request) ([]byte, error) {
	cm := req.Params()
	return json.Marshal(cm)
}

func (cc *CodeC) Decompress(raw []byte) ([]byte, error) {
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
	return msg, err
}
