package ftx

import (
	"context"
	"net/http"
	"testing"
	"time"

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

	httpmock.RegisterResponder(http.MethodGet, "https://ftx.com/api/futures", httpmock.NewStringResponder(200, data))
	ctx := context.Background()
	client := NewRestClient("", "")
	if err := client.Init(ctx); err != nil {
		t.Fatalf("load context fail %s", err.Error())
	}

	if _, err := client.ParseFutureSymbol("BTC-0229"); err == nil {
		t.Errorf("expect BTC-0229 error")
	}

	if f, _ := client.ParseFutureSymbol("BTC-0329"); f.String() != "BTC-0329" ||
		!f.SettleTime().Equal(time.Date(2019, 3, 29, 3, 0, 0, 0, time.UTC)) {
		t.Errorf("bad future symbol %v %v", f, f.SettleTime())
	}

	if s, _ := client.ParseSwapSymbol("BTC-PERP"); s.String() != "BTC-PERP" {
		t.Errorf("bad swap symbol %v", s)
	}
}
