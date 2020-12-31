package deribit

import (
	"encoding/json"
	"testing"
)

func TestOrderBookEndPoint(t *testing.T) {
	raw := `{
    "underlying_price": 29153.69,
    "underlying_index": "SYN.BTC-15JAN21",
    "timestamp": 1609384730848,
    "stats": {
      "volume": 159.6,
      "price_change": 11.9497,
      "low": 0.0675,
      "high": 0.101
    },
    "state": "open",
    "settlement_price": 0.08,
    "open_interest": 303.4,
    "min_price": 0.0525,
    "max_price": 0.1375,
    "mark_price": 0.09096156,
    "mark_iv": 87.2,
    "last_price": 0.089,
    "interest_rate": 0,
    "instrument_name": "BTC-15JAN21-28000-C",
    "index_price": 28930.14,
    "greeks": {
      "vega": 22.57564,
      "theta": -64.77785,
      "rho": 6.46787,
      "gamma": 0.00007,
      "delta": 0.62396
    },
    "estimated_delivery_price": 28930.14,
    "change_id": 24804000251,
    "bids": [
      [
        0.0895,
        23.4
      ],
      [
        0.051,
        0.5
      ],
      [
        0.001,
        1
      ]
    ],
    "bid_iv": 85.2,
    "best_bid_price": 0.0895,
    "best_bid_amount": 23.4,
    "best_ask_price": 0.0955,
    "best_ask_amount": 0.5,
    "asks": [
      [
        0.0955,
        0.5
      ],
      [
        0.108,
        1
      ],
      [
        0.12,
        1
      ],
      [
        0.121,
        2.1
      ],
      [
        0.25,
        12
      ]
    ],
    "ask_iv": 92.94
}`
	var rbd RestBookData
	if err := json.Unmarshal([]byte(raw), &rbd); err != nil {
		t.Fatalf("unamrshal fail error=%s", err.Error())
	}

	ob, err := rbd.Transform(nil)
	if err != nil {
		t.Fatalf("transform fail error=%s", err.Error())
	}

	if ob.Asks[0].Amount != 0.5 || ob.Asks[0].Price != 0.0955 || len(ob.Asks) != 5 ||
		ob.Bids[0].Price != 0.0895 || ob.Bids[0].Amount != 23.4 || len(ob.Bids) != 3 {
		t.Errorf("bad asks=%v bid=%v", ob.Asks, ob.Bids)
	}
}
