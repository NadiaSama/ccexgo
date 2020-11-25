package okex

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/okex"
)

func TestAll(t *testing.T) {
	key := os.Getenv("O_KEY")
	secret := os.Getenv("O_SECRET")
	passphrase := os.Getenv("O_PASSPHRASE")
	if key == "" || secret == "" || passphrase == "" {
		t.Fatalf("missing auth message '%s' '%s' '%s'", key, secret, passphrase)
	}

	client := okex.NewRestClient(key, secret, passphrase, "www.okex.com")
	ctx := context.Background()
	sym1 := okex.NewSpotSymbol("btc", "usdt")
	fee, err := client.FeeRate(ctx, sym1)
	if err != nil {
		t.Fatalf("load fee fail %v", err.Error())
	}

	if len(fee) != 1 || fee[0].Symbol.String() != "BTC-USDT" {
		t.Errorf("test feerate fail %v", fee)
	}
	fmt.Printf("%s %f %f\n", fee[0].Symbol.String(), fee[0].Maker, fee[0].Taker)
}
