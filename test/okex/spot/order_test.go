package spot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/exchange/okex/spot"
	"github.com/shopspring/decimal"
)

func TestSpotAll(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := spot.NewTestRestClient("", "", "")
	ch := make(chan interface{}, 4)

	var ticker spot.Ticker

	wsclient := okex.NewTESTWSClient("", "", "", ch)
	if err := wsclient.Run(ctx); err != nil {
		t.Fatalf("websocket run fail error=%s", err.Error())
	}

	if err := spot.Init(ctx, true); err != nil {
		t.Fatalf("init fail error=%s", err.Error())
	}
	sym, err := spot.ParseSymbol("MNBTC-MNUSDT")
	if err != nil {
		t.Errorf("init symbol fail error=%s", err.Error())
	}
	if err := wsclient.Subscribe(ctx, spot.NewTickerChannel(sym)); err != nil {
		t.Fatalf("subscirbe fail")
	}

	var price decimal.Decimal
	select {
	case raw := <-ch:
		notify := raw.(*exchange.WSNotify)
		ts := notify.Data.(*spot.Ticker)

		ticker = *ts
		price = ticker.BestBid

	case <-time.After(time.Second * 5):
		price = decimal.NewFromFloat(20000.0)

	}

	order, err := client.CreateOrder(ctx, &exchange.OrderRequest{
		Symbol: sym,
		Amount: decimal.NewFromFloat(0.01),
		Price:  price.Mul(decimal.NewFromFloat(0.995)),
		Side:   exchange.OrderSideBuy,
		Type:   exchange.OrderTypeMarket,
	})
	if err != nil {
		t.Fatalf("create order fail error=%s", err.Error())
	}

	o, err := client.FetchOrder(ctx, order)
	if err != nil {
		t.Fatalf("fetch order fail errro=%s", err.Error())
	}

	fmt.Printf("id=%s symbol=%s amount=%s price=%s filled=%s status=%v side=%s\n",
		o.ID.String(), o.Symbol.String(), o.Amount.String(), o.AvgPrice.String(), o.Filled.String(), o.Status, o.Side)

	return
}
