package huobi

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

func TestAll(t *testing.T) {
	key := os.Getenv("H_KEY")
	secret := os.Getenv("H_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		t.Fatalf("missing H_KEY, H_SECRET")
	}
	client := huobi.NewRestClient(key, secret, "api.huobi.pro")
	ctx := context.Background()
	if err := client.Init(ctx); err != nil {
		t.Fatalf("load pairs fail %v", err)
	}

	rates, err := client.FeeRate(ctx, client.NewSpotSymbol("BTC", "USDT"),
		client.NewSpotSymbol("ETH", "USDT"), client.NewSpotSymbol("LTC", "usdt"))

	if err != nil {
		t.Errorf("load rate fail %v", err.Error())
	}

	fmt.Printf("GOT RATE %v\n", rates)

	fs := client.GetFutureSymbols("BTC")

	m := map[string]string{}
	for _, f := range fs {
		m[f.WSSub()] = f.String()
	}

	data := make(chan interface{}, 8)
	ws := huobi.NewWSClient("wss://api.hbdm.com/ws", data, m)
	if err := ws.Run(ctx); err != nil {
		t.Fatalf("run loop fail %s", err.Error())
	}
	if err := ws.Subscribe(ctx, exchange.SubTypeTrade, fs[0], fs[1], fs[2], fs[3]); err != nil {
		t.Fatalf("subscribe fail %s", err.Error())
	}

	statis := map[string]int{}
	for _, f := range fs {
		statis[f.WSSub()] = 0
	}

	timer := time.NewTimer(time.Second * 3)
	done := time.NewTimer(time.Second * 5)
	for {
		select {
		case d := <-data:
			nt := d.(*huobi.NotifyTrade)
			//market.BTC20200925.BTC_CQ.trade.detail
			fields := strings.Split(nt.Chan, ".")
			statis[fields[2]]++

		case <-timer.C:
			for k, v := range statis {
				if v == 0 {
					t.Errorf("no %s data", k)
				}
			}
			ws.Close()

		case <-ws.Done():
			return

		case <-done.C:
			t.Fatalf("time out!")
		}
	}
}
