package jex

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/jex"
)

func TestRequest(t *testing.T) {
	key := os.Getenv("J_KEY")
	secret := os.Getenv("J_SECRET")
	ctx := context.Background()
	client := jex.NewClient(ctx, key, secret)

	_, err := client.Request("GET", "/api/v1/optionInfo", nil, false)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	params := map[string]string{
		"symbol":   "BTC0425CALL",
		"price":    "0.1",
		"quantity": "1.0",
		"side":     "buy",
		"type":     "LIMIT",
	}
	raw, err := client.Request("POST", "/api/v1/option/order", params, true)
	if err != nil {
		t.Errorf("put order fail %v", err)
	}
	fmt.Printf("%s\n", string(raw))
}
