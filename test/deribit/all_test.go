package deribit

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/deribit"
)

func TestAll(t *testing.T) {
	baseCtx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	key := os.Getenv("D_KEY")
	secret := os.Getenv("D_SECRET")
	if key == "" || secret == "" {
		t.Fatalf("missing env D_KEY D_SECRET")
	}

	client := deribit.NewClient(key, secret, time.Second*5, true)
	if err := client.Run(baseCtx); err != nil {
		t.Fatalf("running the loop fail %s", err.Error())
	}

	spot, _ := deribit.ParseSpotSymbol("btc_usd")
	if err := client.Subscribe(baseCtx, exchange.SubTypeIndex, spot); err != nil {
		t.Fatalf("subscribe index fail %+v %v %s", err, reflect.TypeOf(err), err.Error())
	}
	instruments, err := client.OptionFetchInstruments(baseCtx, "BTC")
	if err != nil {
		t.Fatalf("load instrument error %s", err.Error())
	}
	index, _ := client.Index(spot)
	fmt.Printf("GOT INDEX %v\n", *index)
	var sym exchange.OptionSymbol
	for _, i := range instruments {
		if i.SettlementPeriod != "day" {
			continue
		}

		if i.Strike > index.Price {
			sym, _ = deribit.ParseOptionSymbol(i.InstrumentName)
			fmt.Printf("GOT SYMBOL %v %v\n", sym, i)
			break
		}

	}

	if err := client.Subscribe(baseCtx, exchange.SubTypeOrderBook, sym); err != nil {
		t.Fatalf("subscribe orderbook fail %s", err.Error())
	}
	//wait goroutine handle orderbook update
	time.Sleep(2 * time.Second)
	orderbook, err := client.OrderBook(sym)
	if err != nil {
		t.Fatalf("load order book fail %s", err.Error())
	}
	fmt.Printf("GOT ORDERBOOK %v\n", *orderbook)
	var price float64
	if len(orderbook.Asks) != 0 {
		price = orderbook.Asks[0].Price - 0.0005
		if price < 0.0 {
			price = -1.0 * price
		}
	} else {
		price = 0.0005
	}

	fmt.Printf("ORDER %s %f\n", sym.String(), price)
	req := exchange.OrderRequest{
		Symbol: sym,
		Price:  price,
		Amount: 0.1,
		Type:   exchange.OrderTypeLimit,
		Side:   exchange.OrderSideBuy,
	}
	//create a order with price will not being executed
	order, err := client.OptionCreateOrder(baseCtx, &req)
	if err != nil {
		t.Fatalf("create order fail %v", err.Error())
	}
	if order.Status != exchange.OrderStatusOpen || !order.Created.Equal(order.Updated) || order.Symbol.String() != sym.String() {
		t.Errorf("bad order status %v", *order)
	}

	if _, err := client.OptionCancelOrder(baseCtx, order); err != nil {
		t.Errorf("cancel order fail %s", err.Error())
	}

	if order, err := client.OptionFetchOrder(baseCtx, order); err != nil {
		t.Errorf("fetch order fail %s", err.Error())
	} else {
		if order.Status != exchange.OrderStatusCancel {
			t.Errorf("test cancel fail %v", *order)
		}
	}

	//test creat a fok order
	if order, err = client.OptionCreateOrder(baseCtx, &req,
		exchange.NewTimeInForceOption(exchange.TimeInForceFOK),
		exchange.NewPostOnlyOption(false),
	); err != nil {
		t.Errorf("test create fok order fail %s", err.Error())
	} else if order.Status != exchange.OrderStatusCancel {
		t.Errorf("fok order executed %v", *order)
	}

	if err := client.UnSubscribe(baseCtx, exchange.SubTypeOrderBook, sym); err != nil {
		t.Errorf("unsubscribe orderbook fail %s", err.Error())
	}
	if err := client.UnSubscribe(baseCtx, exchange.SubTypeIndex, spot); err != nil {
		t.Errorf("unsubscribe index fail %s", err.Error())
	}
}
