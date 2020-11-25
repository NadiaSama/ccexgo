package okex

import (
	"context"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
	}
)

const (
	OkexWSAddr     = "wss://real.okex.com:8443/ws/v3"
	OkexTESTWSAddr = "wss://real.okex.com:8443/ws/v3?brokerId=9999"
	OKEX           = "okex"
	opSubscribe    = "subscribe"
	opUnSubscribe  = "unsubscribe"
)

func NewWSClient(data chan interface{}) *WSClient {
	return newWSClient(OkexWSAddr, data)
}

//NewTESTWSClient return a wsclient for okex testnet
func NewTESTWSClient(data chan interface{}) *WSClient {
	return newWSClient(OkexTESTWSAddr, data)
}

func newWSClient(addr string, data chan interface{}) *WSClient {
	ret := &WSClient{
		data: data,
	}
	codec := NewCodeC()
	ret.WSClient = exchange.NewWSClient(addr, codec, ret)
	return ret
}

//Subscribe in order to get subscribe result only one channle can subscribe each time
func (ws *WSClient) Subscribe(ctx context.Context, channel ...exchange.Channel) error {
	if len(channel) != 1 {
		return errors.Errorf("only one channel can subscribe each time")
	}

	arg := channel[0]
	cm := callParam{
		OP:   opSubscribe,
		Args: []string{arg.String()},
	}

	var r response
	if err := ws.Call(ctx, opSubscribe, opSubscribe, &cm, &r); err != nil {
		return errors.WithMessagef(err, "subscribe error '%s'", arg.String())
	}
	return nil
}

//Run start the websocket loop and create a goroutine which
//will send ping message to okex server periodically
func (ws *WSClient) Run(ctx context.Context) error {
	if err := ws.WSClient.Run(ctx); err != nil {
		return err
	}

	//period send ping message check the ws conn is correct
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				var msg map[string]interface{}
				if err := ws.Call(ctx, idPingPong, pingMsg, pingMessage, &msg); err != nil {
					//TODO make rpc.Conn fail method public?
					ws.WSClient.Close()
					return
				}
			}
		}
	}()
	return nil
}

func (ws *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	data := &exchange.WSNotify{
		Exchange: OKEX,
		Chan:     notify.Method,
		Data:     notify.Params,
	}
	select {
	case ws.data <- data:
	default:
		return
	}
}
