package okex

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*exchange.WSClient
		data       chan interface{}
		Key        string
		Secret     string
		PassPhrase string
	}
)

const (
	OkexWSAddr     = "wss://real.okex.com:8443/ws/v3"
	OkexTESTWSAddr = "wss://real.okex.com:8443/ws/v3?brokerId=9999"
	OKEX           = "okex"
	opSubscribe    = "subscribe"
	opUnSubscribe  = "unsubscribe"
	opLogin        = "login"
)

func NewWSClient(key, secret, passPhrase string, data chan interface{}) *WSClient {
	return newWSClient(OkexWSAddr, key, secret, passPhrase, data)
}

//NewTESTWSClient return a wsclient for okex testnet
func NewTESTWSClient(key, secret, passPhrase string, data chan interface{}) *WSClient {
	return newWSClient(OkexTESTWSAddr, key, secret, passPhrase, data)
}

func newWSClient(addr, key, secret, passPhrase string, data chan interface{}) *WSClient {
	ret := &WSClient{
		data:       data,
		Key:        key,
		Secret:     secret,
		PassPhrase: passPhrase,
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

			case <-ws.Done():
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

func (ws *WSClient) Auth(ctx context.Context) error {
	timestamp := strconv.FormatFloat(float64(time.Now().UnixNano()/1e6/1000), 'f', -1, 64)
	h := hmac.New(sha256.New, []byte(ws.Secret))
	h.Write([]byte(timestamp + "GET/users/self/verify"))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	cm := callParam{
		OP:   opLogin,
		Args: []string{ws.Key, ws.PassPhrase, timestamp, sign},
	}

	var msg map[string]interface{}
	if err := ws.Call(ctx, opLogin, opLogin, &cm, &msg); err != nil {
		return errors.WithMessage(err, "okex login error")
	}
	return nil
}
