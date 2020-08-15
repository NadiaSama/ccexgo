package ftx

import (
	"testing"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/shopspring/decimal"
)

func TestMessageDecode(t *testing.T) {
	cc := NewCodeC(map[string]exchange.Symbol{
		"ADA-PERP": newSwapSymbol("ADA"),
	})

	e := []byte(`{"channel": "", "market": "", "type": "error", "code": 1001, "msg": "not login"}`)
	if resp, err := cc.Decode(e); err != nil {
		t.Errorf("parse error fail %s", err.Error())
	} else {
		if r := resp.(*rpc.Result); r.Error == nil {
			t.Errorf("expect error fail %v", *r)
		}
	}

	s := []byte(`{"channel": "ch", "market": "m", "type": "subscribed"}`)
	if resp, err := cc.Decode(s); err != nil {
		t.Errorf("parse error fail %s", err.Error())
	} else {
		if r := resp.(*rpc.Result); r.ID != "chm" {
			t.Errorf("bad id %v", *r)
		}
	}

	i := []byte(`{"type": "info", "code": 20001}`)
	if _, err := cc.Decode(i); err == nil {
		t.Errorf("parse info fail")
	} else {
		if _, ok := err.(*rpc.StreamError); !ok {
			t.Errorf("expect streamerror %v", err)
		}
	}

	o := []byte(`{"channel": "orders", "type": "update", "data": {"id": 7670152177, "clientId": null, "market": "ADA-PERP", "type": "limit", "side": "buy", "price": 0.12, "size": 1.0, "status": "new", "filledSize": 0.0, "remainingSize": 1.0, "reduceOnly": false, "liquidation": false, "avgFillPrice": null, "postOnly": false, "ioc": false, "createdAt": "2020-08-15T03:51:55.819920+00:00"}}`)
	if notify, err := cc.Decode(o); err != nil {
		t.Errorf("parse info fail %s", err.Error())
	} else {
		n := notify.(*rpc.Notify)
		order := n.Params.(*exchange.Order)

		if order.ID.String() != "7670152177" || !order.AvgPrice.Equal(decimal.Zero) || order.Symbol.String() != "ADA-PERP" ||
			order.Side != exchange.OrderSideBuy || order.Status != exchange.OrderStatusOpen {
			t.Errorf("bad order %v", *order)
		}
	}
	o2 := []byte(`{"channel": "orders", "type": "update", "data": {"id": 7670152177, "clientId": null, "market": "ADA-PERP", "type": "limit", "side": "buy", "price": 0.12, "size": 1.0, "status": "new", "filledSize": 0.0, "remainingSize": 1.0, "reduceOnly": false, "liquidation": false, "avgFillPrice": 1.02, "postOnly": false, "ioc": false, "createdAt": "2020-08-15T03:51:55.819920+00:00"}}`)
	if notify, err := cc.Decode(o2); err != nil {
		t.Errorf("parse info fail %s", err.Error())
	} else {
		n := notify.(*rpc.Notify)
		order := n.Params.(*exchange.Order)

		if !order.AvgPrice.Equal(decimal.NewFromFloat(1.02)) {
			t.Errorf("bad order %v", *order)
		}
	}
}
