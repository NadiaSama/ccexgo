package deribit

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/deribit"
)

func TestAll(t *testing.T) {
	baseCtx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	key := os.Getenv("D_KEY")
	secret := os.Getenv("D_SECRET")
	if key == "" || secret == "" {
		t.Fatalf("missing env D_KEY D_SECRET")
	}

	client := deribit.NewClient(key, secret, true)
	channels := []string{}
	if err := client.Run(baseCtx); err != nil {
		t.Fatalf("running the loop fail %s", err.Error())
	}

	if err := client.Subscribe(baseCtx, "deribit_price_index.btc_usd"); err != nil {
		t.Fatalf("subscribe index fail %s", err.Error())
	}
	channels = append(channels, "deribit_price_index.btc_usd")
	instruments, err := client.OptionFetchInstruments(baseCtx, "BTC")
	if err != nil {
		t.Fatalf("load instrument error %s", err.Error())
	}
	spot, _ := deribit.ParseSpotSymbol("btc_usd")
	index, _ := client.Index(spot)
	fmt.Printf("GOT INDEX %v\n", *index)
	var sym exchange.OptionSymbol
	for _, i := range instruments {
		if i.SettlementPeriod != "day" {
			continue
		}
		if i.Strike > index.Price {
			sym, _ = deribit.PraseOptionSymbol(i.InstrumentName)
			break
		}
	}

	if err := client.Subscribe(baseCtx, fmt.Sprintf("book.%s.raw", sym.String())); err != nil {
		t.Fatalf("subscribe orderbook fail %s", err.Error())
	}
	channels = append(channels, fmt.Sprintf("book.%s.raw", sym.String()))
	//wait goroutine handle orderbook update
	time.Sleep(100 * time.Millisecond)
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
	order, err := client.OptionCreateOrder(baseCtx, &req)
	if err != nil {
		t.Fatalf("create order fail %v", err.Error())
	}
	if order.Status != exchange.OrderStatusOpen || !order.Created.Equal(order.Updated) {
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

	if err := client.UnSubscribe(baseCtx, channels...); err != nil {
		t.Errorf("unsubscribe fail %s", err.Error())
	}
}
