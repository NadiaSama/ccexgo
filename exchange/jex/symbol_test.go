package jex

import (
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

func TestSymbol(t *testing.T) {
	sym := NewOptionSymbol("EOS", 100.0, time.Date(2020, 3, 4, 0, 0, 0, 0, time.UTC), exchange.OptionTypePut)
	SetSymbol(sym)

	if _, err := ParseSymbol("EOS0304PUT"); err != nil {
		t.Errorf("load symbol fail %v", err)
	}
}
