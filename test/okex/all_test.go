package okex

import (
	"context"
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
	sym1 := client.NewSpotSymbol("btc", "usdt")
	sym2 := client.NewSpotSymbol("eth", "usdt")
	fee, err := client.FeeRate(ctx, sym1, sym2)
	if err != nil {
		t.Fatalf("load fee fail %v", err.Error())
	}

	if len(fee) != 2 || fee[0].Symbol.String() != "BTC-USDT" || fee[1].Symbol.String() != "ETH-USDT" ||
		fee[1].Maker != fee[0].Maker || fee[1].Taker != fee[0].Taker {
		t.Errorf("test feerate fail %v", fee)
	}
}
