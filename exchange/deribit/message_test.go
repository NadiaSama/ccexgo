package deribit

import (
	"testing"

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
