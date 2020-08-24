package spot

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
                "12"
            ],
            [
                "11645.9",
                "0.51397233",
                "2"
            ],
            [
                "11646.4",
                "0.001",
                "1"
            ],
            [
                "11646.5",
                "0.112",
                "1"
            ],
            [
                "11646.6",
                "0.065",
                "2"
            ]
        ],
        "bids": [
            [
                "11645.7",
                "8.82299797",
                "18"
            ],
            [
                "11645.6",
                "0.289",
                "1"
            ],
            [
                "11645.5",
                "0.274",
                "2"
            ],
            [
                "11645.3",
                "0.001",
                "1"
            ],
            [
                "11645.2",
                "0.19998001",
                "1"
            ]
        ],
        "instrument_id": "BTC-USDT",
        "timestamp": "2020-08-24T02:58:24.064Z"
    }
]`
	p, err := parseDepth5(table, action, []byte(data))
	if err != nil {
		t.Fatalf("parse fail '%s'", err.Error())
	}

	d5 := p.Params.(*Depth5)
	o, _ := d5.Bids[0].Orders.Float64()
	if d5.Symbol.String() != "BTC-USDT" || d5.Bids[0].Price.String() != "11645.7" || d5.Bids[0].Amount.String() != "8.82299797" || int(o) != int(18) {
		t.Fatalf("compare fail %s %v %v %s %s", d5.Symbol.String(), d5.Bids[0].Price.String() != "11645.7",
			d5.Bids[0].Amount.String() != "8.82299797", d5.Bids[0].Price.String(), d5.Bids[0].Amount.String())
	}
}
