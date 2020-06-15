package deribit

import (
	"encoding/json"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

func TestDecode(t *testing.T) {
	cc := &Codec{}
	message := `{
		"jsonrpc": "2.0",
		"id": 8163,
		"error": {
			"code": 11050,
			"message": "bad_request"
		},
		"testnet": false,
		"usIn": 1535037392434763,
		"usOut": 1535037392448119,
		"usDiff": 13356
	}`
	resp, err := cc.Decode([]byte(message))
	if err != nil {
		t.Fatalf("decode fail %s", err.Error())
	}
	result := resp.(*rpc.Result)
	if result.ID.Num != 8163 || result.Error.Code != 11050 || result.Error.Message != "bad_request" {
		t.Errorf("bad result %v", *result)
	}
}

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
	notify, _ := parseNotifyBook(json.RawMessage(raw), "test")
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
	if _, err := parseNotifyBook(json.RawMessage(bad), ""); err == nil {
		t.Errorf("test bad format fail")
	}
}
