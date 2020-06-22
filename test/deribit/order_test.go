package deribit

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/deribit"
)

func TestOrderBuy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := os.Getenv("D_KEY")
	secret := os.Getenv("D_SECRET")
	if key == "" || secret == "" {
		t.Fatalf("missing env D_KEY D_SECRET")
	}

	client := deribit.NewClient(key, secret, true)
	if err := client.Run(ctx); err != nil {
		t.Fatalf("running the loop fail %s", err.Error())
	}
	sym, _ := deribit.PraseOptionSymbol("BTC-26JUN20-9000-C")
	order, err := client.OptionCreateOrder(ctx, sym, exchange.OrderSideBuy, 0.001, 0.1, exchange.OrderTypeLimit)
	if err != nil {
		t.Fatalf("create order fail %v", err.Error())
	}
	if order.Status != exchange.OrderStatusOpen || !order.Created.Equal(order.Updated) {
		t.Errorf("bad order status %v", *order)
	}

	if _, err := client.OptionCancelOrder(ctx, order); err != nil {
		t.Errorf("cancel order fail %s", err.Error())
	}

	if order, err := client.OptionFetchOrder(ctx, order); err != nil {
		t.Errorf("fetch order fail %s", err.Error())
	} else {
		if order.Status != exchange.OrderStatusCancel {
			t.Errorf("test cancel fail %v", *order)
		}
	}
}
