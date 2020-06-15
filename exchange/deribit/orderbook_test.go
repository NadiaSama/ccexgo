package deribit

import (
	"encoding/json"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange"
)

func TestNotifyBookTest(t *testing.T) {
	raw := `{
      "type" : "change",
      "timestamp" : 1554373911330,
      "prev_change_id" : 297217,
      "instrument_name" : "BTC-PERPETUAL",
      "change_id" : 297218,
      "bids" : [
        [
          "new",
          5041.94,
          0
        ],
        [
          "new",
          5042.34,
          10
        ]
      ],
      "asks" : [
		  [
			  "change",
			  5041.2,
			  12
		  ]
      ]
	}`
	notify, _ := parseNotifyBook(&Notify{Data: json.RawMessage(raw), Channel: "1.2"})
	n := notify.Params.(*exchange.OrderBookNotify)
	if n.Asks[0].Amount != 12 || n.Asks[0].Price != 5041.2 || n.Bids[0].Price != 5041.94 ||
		n.Bids[0].Amount != 0 || n.Bids[1].Price != 5042.34 || n.Bids[1].Amount != 10 {
		t.Errorf("bad notify %v", *n)
	}

	bad := `{
      "type" : "change",
      "timestamp" : 1554373911330,
      "instrument_name" : "BTC-PERPETUAL",
      "bids" : [
        [
          "new",
          "5042.34",
          10
        ]
      ],
      "asks" : [ ]}`
	if _, err := parseNotifyBook(&Notify{Data: json.RawMessage(bad), Channel: "1.2"}); err == nil {
		t.Errorf("test bad format fail")
	}
}
