package huobi

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

func TestTrades(t *testing.T) {
	raw := `{
    "ch": "market.BTC_NW.trade.detail",
    "ts": 1539831709042,
    "tick": {
        "id": 265842227,
        "ts": 1539831709001,
        "data": [{
            "amount": 20,
            "ts": 1539831709001,
            "id": 265842227259096443,
            "price": 6742.25,
            "direction": "buy"
        }]
    }
}`
	var buf bytes.Buffer
	b := gzip.NewWriter(&buf)
	b.Write([]byte(raw))
	b.Close()
	cc := CodeC{
		codeMap: map[string]string{
			"BTC_NW": "btc1234",
		},
	}

	if resp, err := cc.Decode(buf.Bytes()); err != nil {
		t.Errorf("parse error %s", err.Error())
	} else {
		notify := resp.(*rpc.Notify)
		if notify.Method != "market.btc1234.trade.detail" {
			t.Errorf("bad method %s", notify.Method)
		}

		if trade := notify.Params.([]Trade); trade[0].Amount != 20.0 || trade[0].Direction != "buy" ||
			trade[0].Price != 6742.25 || trade[0].TS != 1539831709001 {
			t.Errorf("bad value %v", trade)
		}
	}
}
