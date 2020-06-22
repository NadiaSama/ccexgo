package exchange

import (
	"context"
	"reflect"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	handlerMsg interface {
		Key() string
	}

	handlerMsgCB func(ds interface{}, msg handlerMsg) interface{}
)

var (
	subTyp2CB = map[reflect.Type]handlerMsgCB{}
)

func subRegister(typ reflect.Type, cb handlerMsgCB) {
	subTyp2CB[typ] = cb
}

//Handler handle notify message
func (c *Client) Handle(_ context.Context, notify *rpc.Notify) {
	c.SubMu.Lock()
	defer c.SubMu.Unlock()

	msg := notify.Params.(handlerMsg)
	typ := reflect.TypeOf(msg)
	cb := subTyp2CB[typ]

	val := c.Sub[msg.Key()]
	val = cb(val, msg)
	c.Sub[msg.Key()] = val
}
