package wallet

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/ftx"
)

func TestBalances(t *testing.T) {
	key := os.Getenv("F_KEY")
	secret := os.Getenv("F_SECRET")
	if key == "" || secret == "" {
		t.Fatalf("missing F_KEY F_SECRET")
	}

	ctx := context.Background()
	client := ftx.NewRestClient(key, secret)
	balance, err := client.Balances(ctx)
	if err != nil {
		t.Fatalf("fetch balance fail %s", err.Error())
	}

	fmt.Printf("GOT BALANCES %v", balance)
}
