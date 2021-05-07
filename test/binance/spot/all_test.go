package spot

import (
	"context"
	"fmt"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/binance/spot"
)

func TestAll(t *testing.T) {
	ctx := context.Background()
	if err := spot.Init(ctx); err != nil {
		t.Fatalf("init fail error=%s", err.Error())
	}

	symbol, _ := spot.ParseSpotSymbol("BTCUSDT")
	s, _ := symbol.(*spot.SpotSymbol)
	fmt.Printf("%+v %s %s\n", *s, s.Base(), s.Quote())

	s2, _ := spot.ParseSpotSymbol("ETHUSDT")
	s, _ = s2.(*spot.SpotSymbol)
	fmt.Printf("%+v\n", *s)
}
