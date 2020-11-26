package swap

import (
	"context"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/exchange/okex/swap"
	"github.com/shopspring/decimal"
)

func TestSwapOrder(t *testing.T) {
	key := "3c81dfc1-1c12-4eae-a3af-84489e45456c"
	secret := "FE87BF6D6A54B4754D20D7488B2EC727"
	passPhrass := ""
	data := make(chan interface{}, 4)
	ctx := context.Background()

	ws := okex.NewTESTWSClient(key, secret, "", data)
	if err := ws.Run(ctx); err != nil {
		t.Fatalf("run websocket fail error=%s", err.Error())
	}
	if err := ws.Auth(ctx); err == nil {
		t.Errorf("test login error fail")
	}

	select {
	case <-ws.Done():
		break
	case <-time.After(time.Second):
		t.Errorf("auto close connection timeout!")
	}

	client := swap.RestClient{
		okex.NewTESTRestClient(key, secret, passPhrass),
	}
	ws = okex.NewTESTWSClient(key, secret, passPhrass, data)
	if err := ws.Run(ctx); err != nil {
		t.Errorf("run fail error=%s", err.Error())
	}
	if err := ws.Auth(ctx); err != nil {
		t.Fatalf("auth fail error=%s", err.Error())
	}

	symbol := okex.NewSwapSymbol("MNBTC-USDT")

	if err := ws.Subscribe(ctx, swap.NewOrderChannel(symbol)); err != nil {
		t.Fatalf("subscribe order fail error=%s", err.Error())
	}

	order, err := client.CreateOrder(ctx, &exchange.OrderRequest{
		Symbol: symbol,
		Side:   exchange.OrderSideBuy,
		Amount: decimal.NewFromFloat(1),
		Price:  decimal.NewFromFloat(17000.0),
		Type:   exchange.OrderTypeLimit,
	})
	if err != nil {
		t.Fatalf("create order fail error=%s", err.Error())
	}

	if err := client.CancelOrder(ctx, order); err != nil {
		t.Fatalf("cancel order fail error=%s", err.Error())
	}

	if _, err := client.CreateOrder(ctx, &exchange.OrderRequest{
		Symbol: symbol,
		Side:   exchange.OrderSideBuy,
		Amount: decimal.NewFromFloat(1),
		Price:  decimal.NewFromFloat(17300.0),
		Type:   exchange.OrderTypeMarket,
	}); err != nil {
		t.Fatalf("create order fail error=%s", err.Error())
	}

	count := 0
	for raw := range data {
		notify := raw.(*exchange.WSNotify)
		orders := notify.Data.([]*exchange.Order)

		t.Logf("got order %+v", *(orders[0]))
		count += 1

		if count == 4 {
			break
		}
	}

	if err := ws.Close(); err != nil {
		t.Errorf("close connection fail error=%s", err.Error())
	}

}
