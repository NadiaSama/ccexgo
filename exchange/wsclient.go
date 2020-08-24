package exchange

import (
	"context"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	WSClient struct {
		rpc.Conn
		handler rpc.Handler
		codec   rpc.Codec
		addr    string
	}

	WSNotify struct {
		Exchange string
		Chan     string
		Data     interface{}
	}

	//Channel a subscribe channel
	Channel interface {
		String() string
	}
)

func NewWSClient(addr string, codec rpc.Codec, handler rpc.Handler) *WSClient {
	return &WSClient{
		addr:    addr,
		codec:   codec,
		handler: handler,
	}
}

func (wc *WSClient) Run(ctx context.Context) error {
	stream, err := rpc.NewWebsocketStream(wc.addr, wc.codec)
	if err != nil {
		return err
	}

	conn := rpc.NewConn(stream)
	wc.Conn = conn
	go wc.Conn.Run(ctx, wc.handler)
	return nil
}
