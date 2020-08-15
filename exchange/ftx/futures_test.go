package ftx

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func TestFuture(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()
	response := `{
  "success": true,
  "result": {
    "ask": 4196,
    "bid": 4114.25,
    "change1h": 0,
    "change24h": 0,
    "description": "Bitcoin March 2019 Futures",
    "enabled": true,
    "expired": false,
    "expiry": "2019-03-29T03:00:00+00:00",
    "index": 3919.58841011,
    "last": 4196,
    "lowerBound": 3663.75,
    "mark": 3854.75,
    "name": "BTC-0329",
    "perpetual": false,
    "postOnly": false,
    "priceIncrement": 0.25,
    "sizeIncrement": 0.001,
    "underlying": "BTC",
    "upperBound": 4112.2,
    "type": "future"
  }
}`

	httpmock.RegisterResponder(http.MethodGet, "https://ftx.com/api/futures/BTC-0329",
		httpmock.NewBytesResponder(200, []byte(response)))

	ctx := context.Background()
	client := NewRestClient("", "")
	future, err := client.Future(ctx, "BTC-0329")
	if err != nil {
		t.Fatalf("get future fail %s", err.Error())
	}

	index := decimal.NewFromFloat(future.Index)
	if future.Ask != 4196 || future.Bid != 4114.25 || !index.Equal(decimal.NewFromFloat(3919.58841011)) ||
		future.Enabled != true {
		t.Errorf("bad future info %v", *future)
	}
}
