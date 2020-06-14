package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexws/exchange/deribit"
	"github.com/NadiaSama/ccexws/internal/rpc"
)

type testHandler struct {
}

func (th *testHandler) Handle(_ context.Context, msg *rpc.Notify) {
	raw, _ := json.Marshal(msg)
	fmt.Printf("GOT %v\n", string(raw))
}

func main() {
	stream, err := rpc.NewWebsocketStream(deribit.WSAddr, &deribit.Codec{})
	if err != nil {
		panic(err)
	}
	conn := rpc.NewConn(stream)
	ctx := context.Background()
	go conn.Run(ctx, &testHandler{})

	var r rpc.Result
	if err := conn.Call(ctx, "public/subscribe", map[string]interface{}{
		"channels": []string{"ticker.BTC-PERPETUAL.100ms"}}, &r); err != nil {
		panic(err)
	}

	fmt.Printf("SUBSCRIBE %v\n", r)

	forever := make(chan struct{})
	<-forever
}
