package ftx

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/jarcoal/httpmock"
)

func TestSymbols(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	data := `{
  "success": true,
  "result": [
    {
      "enabled": true,
      "expired": false,
      "expiry": "2019-03-29T03:00:00+00:00",
      "name": "BTC-0329",
      "perpetual": false,
      "underlying": "BTC",
      "type": "future"
	},
    {
      "enabled": true,
      "expired": false,
      "expiry": "2019-03-29T03:00:00+00:00",
      "name": "BTC-PERP",
      "perpetual": false,
      "underlying": "BTC",
      "type": "perpetual"
	},
    {
      "enabled": true,
      "expired": true,
      "expiry": "2019-03-29T03:00:00+00:00",
      "name": "BTC-0229",
      "perpetual": false,
      "underlying": "BTC",
      "type": "perpetual"
	}
  ]
}`

	mdata := `{"success": true, "result": []}`

	httpmock.RegisterResponder(http.MethodGet, "https://ftx.com/api/futures", httpmock.NewStringResponder(200, data))
	httpmock.RegisterResponder(http.MethodGet, "https://ftx.com/api/markets", httpmock.NewStringResponder(200, mdata))
	ctx := context.Background()
	if err := Init(ctx); err != nil {
		t.Fatalf("load context fail %s", err.Error())
	}

	if _, err := ParseSymbol("BTC-0229"); err == nil {
		t.Errorf("expect BTC-0229 error")
	}

	if r, _ := ParseSymbol("BTC-0329"); r.String() != "BTC-0329" ||
		!r.(exchange.FuturesSymbol).SettleTime().Equal(time.Date(2019, 3, 29, 3, 0, 0, 0, time.UTC)) {
		t.Errorf("bad future symbol %v,  %v,", r.String(), r.(exchange.FuturesSymbol).SettleTime())
	}

	if s, _ := ParseSymbol("BTC-PERP"); s.String() != "BTC-PERP" {
		t.Errorf("bad swap symbol %v", s)
	}
}
