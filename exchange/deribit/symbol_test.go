package deribit

import "testing"

func TestSymbol(t *testing.T) {
	v := "BTC-19JUN20-10250-C"
	sym, err := PraseOptionSymbol(v)
	if err != nil {
		t.Fatalf("parse error %s", err.Error())
	}

	if sym.String() != v {
		t.Errorf("bad symbol %s", sym.String())
	}
}
