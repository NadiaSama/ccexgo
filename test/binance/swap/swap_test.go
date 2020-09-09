package swap

import (
	"context"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/binance/swap"
)

func TestSwapClient(t *testing.T) {
	rc := swap.NewRestClient("", "")
	ctx := context.Background()

	if err := rc.Init(ctx); err != nil {
		t.Fatalf("load symbols fail %s", err.Error())
	}

	symbols := rc.Symbols()
	symbol := symbols["ADAUSDT"]
	if symbol.PricePrecision.String() != "0.00001" || symbol.MinAmount.String() != "1" || symbol.AmountPrecision.String() != "1" ||
		symbol.BaseAsset != "ADA" {
		t.Fatalf("not equal %+v", *symbol)
	}

	symbol = symbols["BCHUSDT"]
	if symbol.AmountPrecision.String() != "0.001" || symbol.MinAmount.String() != "0.001" ||
		symbol.PricePrecision.String() != "0.01" || symbol.BaseAsset != "BCH" {
		t.Fatalf("not equal %+v", *symbol)
	}
}
