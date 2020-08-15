package jex

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/jex"
	"github.com/shopspring/decimal"
)

func TestRequest(t *testing.T) {
	key := os.Getenv("J_KEY")
	secret := os.Getenv("J_SECRET")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client := jex.NewClient(key, secret)

	req := &exchange.OrderRequest{
		Symbol: jex.NewOptionSymbol("EOS", time.Date(2020, 7, 6, 11, 0, 0, 0, time.UTC), 1.234, exchange.OptionTypeCall),
		Price:  decimal.NewFromFloat(0.3),
		Amount: decimal.NewFromFloat(0.1),
		Type:   exchange.OrderTypeLimit,
		Side:   exchange.OrderSideBuy,
	}
	order, err := client.OptionCreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("put order fail %v", err)
	}
	fmt.Printf("new order %v", *order)

	co, err := client.OptionCancelOrder(ctx, order)
	if err != nil {
		t.Errorf("cancel order fail %v", err)
	} else {
		fmt.Printf("cancel result %v\n", *co)
	}
}
