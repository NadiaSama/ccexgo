package future

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange/huobi/future"

	"github.com/NadiaSama/ccexgo/exchange"
)

func TestAll(t *testing.T) {
	key := os.Getenv("H_KEY")
	secret := os.Getenv("H_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		t.Fatalf("missing H_KEY, H_SECRET")
	}
	client := future.NewRestClient(key, secret)
	ctx, cancel := context.WithCancel(context.Background())
	if err := client.Init(ctx); err != nil {
		t.Fatalf("future client init fail %s", err.Error())
	}
	fs := client.GetFutureSymbols("BTC")

	m := map[string]string{}
	for _, f := range fs {
		m[f.WSSub()] = f.String()
	}

	data := make(chan interface{}, 8)
	ws := future.NewWSClient(m, data)
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
			nt := d.(*exchange.WSNotify)
			//market.BTC20200925.BTC_CQ.trade.detail
			fmt.Printf("%v\n", *nt)
			fields := strings.Split(nt.Chan, ".")
			statis[fields[2]]++

		case <-timer.C:
			for k, v := range statis {
				if v == 0 {
					t.Errorf("no %s data", k)
				}
			}
			cancel()

		case <-ws.Done():
			return

		case <-done.C:
			t.Fatalf("time out!")
		}
	}
}
