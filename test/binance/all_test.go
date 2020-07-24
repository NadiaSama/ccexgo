package binance

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/binance"
)

func TestAll(t *testing.T) {
	key := os.Getenv("B_KEY")
	secret := os.Getenv("B_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		t.Fatalf("missing B_KEY B_SECRET env")
	}

	ctx := context.Background()
	client := binance.NewRestClient(key, secret, "api.binance.com")
	if err := client.Init(ctx); err != nil {
		t.Fatalf("init pair fail %v", err.Error())
	}

	if rate, err := client.FeeRate(ctx, client.NewSpotSymbol("btc", "usdt")); err != nil {
		t.Errorf("got feerate fail %v", err)
	} else {
		fmt.Printf("Got Rate %v", rate)
	}
}
