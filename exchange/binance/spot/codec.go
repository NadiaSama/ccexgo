package spot

import (
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type (
	// CodeC used to decode binance websocket notify message to coresponding struct
	// and encode sbuscribe request
	CodeC struct {
		*binance.CodeC
	}
)

const (
	MethodSubscibe = "SUBSCRIBE"
)

func NewCodeC() *CodeC {
	return &CodeC{
		CodeC: binance.NewCodeC(),
	}
}

// Decode binance websocket notify message
func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	return cc.DecodeByCB(raw, func(g *gjson.Result) (rpc.Response, error) {
		if g.Get("u").Exists() {
			tn := ParseBookTickerNotify(g)
			return &rpc.Notify{Params: tn, Method: "bookTicker"}, nil
		}

		return nil, errors.Errorf("bad notify msg=%s", g.Raw)
	})
}
