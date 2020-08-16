package test_ws

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/ftx"
)

type ()

func TestOrderWS(t *testing.T) {
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

	if err := client.Subscribe(ctx, exchange.SubTypePrivateTrade); err != nil {
		t.Fatalf("subscribe order fail %s", err.Error())
	}

	for {
		select {
		case raw := <-ch:
			notify := raw.(*exchange.WSNotify)
			order := notify.Data.(*ftx.Fill)
			fmt.Printf("GOT NOTIFY %v\n", *order)
		}
	}

}
