package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange/deribit"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type testHandler struct {
}

func (th *testHandler) Handle(_ context.Context, msg *rpc.Notify) {
	raw, _ := json.Marshal(msg)
	fmt.Printf("GOT %v\n", string(raw))
}

func TestOrderBook(t *testing.T) {
	stream, err := rpc.NewWebsocketStream(deribit.WSTestAddr, &deribit.Codec{})
	if err != nil {
		t.Fatalf("create stream error %v", err)
	}
	conn := rpc.NewConn(stream)
	ctx := context.Background()
	go conn.Run(ctx, &testHandler{})

	var result []string
	sub := []string{"book.BTC-PERPETUAL.100ms"}
	if err := conn.Call(ctx, "public/subscribe", map[string]interface{}{
		"channels": []string{"book.BTC-PERPETUAL.100ms"}}, &result); err != nil {
		t.Fatalf("subscribe channel fail %v", err.Error())
	}

	if len(result) != 1 || sub[0] != result[0] {
		t.Fatalf("subscribe fail %v", result)
	}

	time.Sleep(time.Second)
}
