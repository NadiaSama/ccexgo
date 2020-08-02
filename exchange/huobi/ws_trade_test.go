package huobi

import "testing"

func TestParseTrades(t *testing.T) {
	raw := []byte(`{
        "id": 265842227,
        "ts": 1539831709001,
        "data": [{
            "amount": 20,
            "ts": 1539831709001,
            "id": 265842227259096443,
            "price": 6742.25,
            "direction": "buy"
        }]
	}`)

	r, _ := ParseTrades(raw)
	if trade := r; trade[0].Amount != 20.0 || trade[0].Direction != "buy" ||
		trade[0].Price != 6742.25 || trade[0].TS != 1539831709001 {
		t.Errorf("bad value %v", trade)
	}
}