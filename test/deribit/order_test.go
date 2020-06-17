package deribit

import (
	"context"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/deribit"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

func TestOrderBuy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := rpc.NewWebsocketStream(deribit.WSTestAddr, &deribit.Codec{})
	if err != nil {
		t.Fatalf("create stream error %v", err)
	}
	conn := rpc.NewConn(stream)
	go conn.Run(ctx, &testHandler{})

	key := os.Getenv("D_KEY")
	secret := os.Getenv("D_SECRET")
	if key == "" || secret == "" {
		t.Fatalf("missing env D_KEY D_SECRET")
	}

	client := deribit.NewClient(conn, key, secret)
	sym, _ := deribit.PraseOptionSymbol("BTC-19JUN20-10250-C")
	order, err := client.OptionCreateOrder(ctx, sym, exchange.OrderSideBuy, 0.001, 0.1, exchange.OrderTypeLimit)
	if err != nil {
		t.Fatalf("create order fail %v", err.Error())
	}
	if order.Status != exchange.OrderStatusOpen || !order.Created.Equal(order.Updated) {
		t.Errorf("bad order status %v", *order)
	}

	if err := client.OptionCancelOrder(ctx, order.ID); err != nil {
		t.Errorf("cancel order fail %s", err.Error())
	}
}
