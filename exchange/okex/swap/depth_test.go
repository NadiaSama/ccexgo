package swap

import (
	"testing"
)

func TestParseDepth5(t *testing.T) {
	table := "spot/depth5"
	action := ""
	data := `[
    {
        "asks": [
            [
                "11645.8",
                "6.72040919",
                "1",
                "12"
            ],
            [
                "11645.9",
                "0.51397233",
                "0.1",
                "2"
            ],
            [
                "11646.4",
                "0.001",
                "0.0",
                "1"
            ],
            [
                "11646.5",
                "0.112",
                "0.112",
                "1"
            ],
            [
                "11646.6",
                "0.065",
                "1",
                "2"
            ]
        ],
        "bids": [
            [
                "11645.7",
                "8.82299797",
                "41.3",
                "18"
            ],
            [
                "11645.6",
                "0.289",
                "0.01",
                "1"
            ],
            [
                "11645.5",
                "0.274",
                "0.001",
                "2"
            ],
            [
                "11645.3",
                "0.001",
                "0.0001",
                "1"
            ],
            [
                "11645.2",
                "0.19998001",
                "0.0023",
                "1"
            ]
        ],
        "instrument_id": "BTC-USDT-SWAP",
        "timestamp": "2020-08-24T02:58:24.064Z"
    }
]`
	p, err := parseDepth5(table, action, []byte(data))
	if err != nil {
		t.Fatalf("parse fail '%s'", err.Error())
	}

	d5 := p.Params.(*Depth5)
	o, _ := d5.Bids[0].Orders.Float64()
	if d5.Symbol.String() != "BTC-USDT-SWAP" || d5.Bids[0].Price.String() != "11645.7" || d5.Bids[0].Amount.String() != "8.82299797" || int(o) != int(18) ||
		d5.Bids[0].Liquid.String() != "41.3" {
		t.Fatalf("compare fail %s %v %v %s %s", d5.Symbol.String(), d5.Bids[0].Price.String() != "11645.7",
			d5.Bids[0].Amount.String() != "8.82299797", d5.Bids[0].Price.String(), d5.Bids[0].Amount.String())
	}
}
