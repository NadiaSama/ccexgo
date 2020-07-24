package deribit

import "testing"

func TestSymbol(t *testing.T) {
	v := "BTC-19JUN20-10250-C"
	sym, err := parseOptionSymbol(v)
	if err != nil {
		t.Fatalf("parse error %s", err.Error())
	}

	if sym.String() != v {
		t.Errorf("bad symbol %s", sym.String())
	}

	v2 := "BTC-1JUL20-9875-P"
	sym2, err := parseOptionSymbol(v2)
	if sym2.String() != v2 {
		t.Errorf("bad symbol %s", sym2.String())
	}
}
