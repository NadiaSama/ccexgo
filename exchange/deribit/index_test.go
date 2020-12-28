package deribit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/shopspring/decimal"
)

func TestIndex(t *testing.T) {
	now := time.Now()
	msg := fmt.Sprintf(`{"index_name": "abc_cba", "price": 1.23, "timestamp": %d}`, now.Unix())

	n := &Notify{
		Data:    json.RawMessage(msg),
		Channel: "digifinex_price_check.btc_usd",
	}

	notify, err := parseNotifyIndex(n)
	if err != nil {
		t.Errorf("parse message fail %s", err.Error())
	}
	in := notify.Params.(*exchange.IndexNotify)
	if !in.Price.Equal(decimal.NewFromFloat(1.23)) || in.Symbol.String() != "abc_cba" {
		t.Errorf("bad notify %v", *in)
	}
}
