package swap

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi/swap"
)

func TestSwapSub(t *testing.T) {
	client := swap.NewRestClient("", "")
	ctx := context.Background()
	if err := client.Init(ctx); err != nil {
		t.Fatalf("swap rest client init fail %s", err.Error())
	}

	sym, err := client.GetSwapContract("BTC")
	if err != nil {
		t.Fatalf("got contract fail %s", err.Error())
	}

	data := make(chan interface{}, 8)
	ws := swap.NewWSClient(data)
	if err := ws.Run(ctx); err != nil {
		t.Fatalf("run fail %s", err.Error())
	}
	if err := ws.Subscribe(ctx, exchange.SubTypeTrade, sym); err != nil {
		t.Fatalf("subscribe fail %s", err.Error())
	}

	var sec time.Duration = 20
	done := time.NewTimer(time.Second * sec)
	timeout := time.NewTimer(time.Second * (sec + 2))

	for {
		select {
		case d := <-data:
			t := d.(*exchange.WSNotify)
			fmt.Printf("GOT %v\n", *t)

		case <-done.C:
			ws.Close()

		case <-timeout.C:
			t.Fatalf("timeout")

		case <-ws.Done():
			fmt.Printf("DONE\n")
			return
		}
	}
}
