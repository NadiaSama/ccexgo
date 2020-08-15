package test_order

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/ftx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func expectOrder(src chan interface{}) (*exchange.Order, error) {
	timer := time.After(time.Second * 2)
	select {
	case raw := <-src:
		notify, ok := raw.(*exchange.WSNotify)
		if !ok {
			return nil, errors.Errorf("got notify fail %v", raw)
		}
		order, ok := notify.Data.(*exchange.Order)
		if !ok {
			return nil, errors.Errorf("got order fail %v", notify)
		}

		return order, nil

	case <-timer:
		return nil, errors.Errorf("fetch order timeout")
	}
}

func TestOrder(t *testing.T) {
	key := os.Getenv("F_KEY")
	secret := os.Getenv("F_SECRET")
	symbol := os.Getenv("F_SYM")

	if key == "" || secret == "" || symbol == "" {
		t.Fatalf("missing F_KEY F_SECRET F_SYM")
	}

	ch := make(chan interface{}, 4)
	ctx := context.Background()
	rest := ftx.NewRestClient(key, secret)
	if err := rest.Init(ctx); err != nil {
		t.Fatalf("init fail %s", err.Error())
	}
	client := rest.NewWSClient(ch)

	if err := client.Run(ctx); err != nil {
		t.Fatalf("run error %s", err.Error())
	}

	if err := client.Auth(ctx, key, secret); err != nil {
		t.Fatalf("auth fail %s", err.Error())
	}

	if err := client.Subscribe(ctx, exchange.SubTypePrivateOrder); err != nil {
		t.Fatalf("subscribe order fail %s", err.Error())
	}

	sym, err := rest.ParseSymbol(symbol)
	if err != nil {
		t.Fatalf("bad symbol %s", err.Error())
	}
	fmt.Printf("GOT SYMBOl %s\n", sym.String())

	future, err := rest.Future(ctx, sym.String())
	if err != nil {
		t.Fatalf("get future fail %s", err.Error())
	}

	fmt.Printf("GOT FUTURE %s %f %f\n", future.Name, future.Bid, future.Ask)
	req := &exchange.OrderRequest{
		Symbol: sym,
		Price:  decimal.NewFromFloat(future.Bid * 0.8),
		Amount: decimal.NewFromInt(1),
		Side:   exchange.OrderSideBuy,
		Type:   exchange.OrderTypeLimit,
	}
	order, err := rest.OrderNew(ctx, req)
	if err != nil {
		t.Fatalf("create order fail %s", err.Error())
	}
	exo, err := expectOrder(ch)
	if err != nil {
		t.Fatalf("read websocket order fail %s", err.Error())
	}
	if !order.Equal(exo) {
		t.Errorf("order not equal %v %v", *order, *exo)
	}

	g, err := rest.OrderFetch(ctx, order)
	if err != nil {
		t.Fatalf("got order fail %s", err.Error())
	}
	if !g.Equal(order) {
		t.Errorf("order not equal %v %v", *g, *order)
	}

	if err := rest.OrderCancel(ctx, order); err != nil {
		t.Errorf("cancel order fail %s", err.Error())
	}
	if o, err := expectOrder(ch); err != nil {
		t.Errorf("websocket read order fail %s", err.Error())
	} else {
		if o.Status != exchange.OrderStatusCancel {
			t.Errorf("bad order status %v", *o)
		}
	}
}
