package deribit

import (
	"context"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/shopspring/decimal"
)

func TestSymbol(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := initSymbol(ctx, false); err != nil {
		t.Fatalf("init symbol error=%s", err.Error())
	}

	client := NewWSClient("", "", nil)
	if err := client.Run(ctx); err != nil {
		t.Fatalf("run client fail error=%s", err.Error())
	}

	inst, err := client.OptionFetchInstruments(ctx, Currencies[0])
	if err != nil {
		t.Fatalf("load instrument fail error=%s", err.Error())
	}

	for _, i := range inst {
		sym, err := ParseOptionSymbol(i.InstrumentName)
		if err != nil {
			t.Fatalf("parse instrument fail err=%s", err.Error())
		}

		if sym.String() != i.InstrumentName {
			t.Fatalf("bad sym string %s", sym.String())
		}
	}

	sym, _ := ParseOptionSymbol(inst[0].InstrumentName)

	amount := decimal.NewFromFloat(1.212345)
	if !exchange.Round(amount, sym.PricePrecision()).Equal(decimal.NewFromFloat(1.212)) {
		t.Errorf("bad round value %s", exchange.Round(amount, sym.PricePrecision()))
	}
}
