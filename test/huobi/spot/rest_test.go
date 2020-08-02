package spot

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/huobi/spot"
)

func TestAll(t *testing.T) {
	key := os.Getenv("H_KEY")
	secret := os.Getenv("H_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		t.Fatalf("missing H_KEY, H_SECRET")
	}
	client := spot.NewRestClient(key, secret)
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
}
