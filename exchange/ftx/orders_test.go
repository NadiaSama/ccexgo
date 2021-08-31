package ftx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

//TestOrdersNew test OrderNew and parseOrder
func TestOrdersNew(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	order := []byte(`{
  "success": true,
  "result": {
    "createdAt": "2019-03-05T09:56:55.728933+00:00",
    "filledSize": 0,
    "future": "XRP-PERP",
    "id": 9596912,
    "market": "XRP-PERP",
    "price": 0.306525,
    "remainingSize": 31431,
    "side": "sell",
    "size": 31431,
    "status": "open",
    "type": "limit",
    "reduceOnly": false,
    "ioc": false,
    "postOnly": false,
    "clientId": null
  }
}`)

	httpmock.RegisterResponder(http.MethodPost, "https://ftx.com/api/orders", func(req *http.Request) (*http.Response, error) {
		bytes, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		var or OrderReq
		if err := json.Unmarshal(bytes, &or); err != nil {
			t.Errorf("bad requests '%s' '%s'", string(bytes), err.Error())
		}

		if or.Market != "XRP-PERP" || or.Price != 1.023 || or.Side != "buy" || or.Size != 10.1234 || or.Type != "market" {
			t.Errorf("bad order param")
		}

		return httpmock.NewBytesResponse(200, order), nil
	})

	ctx := context.Background()
	client := NewRestClient("", "")
	xrpS := newSwapSymbol("XRP")
	symbolMap["XRP-PERP"] = xrpS
	req := exchange.OrderRequest{
		Symbol: xrpS,
		Amount: decimal.NewFromFloat(10.1234),
		Price:  decimal.NewFromFloat(1.023),
		Side:   exchange.OrderSideBuy,
		Type:   exchange.OrderTypeStopLimit,
	}

	if _, err := client.OrderNew(ctx, &req); err == nil {
		t.Errorf("test bad order type fali")
	}

	req.Type = exchange.OrderTypeMarket
	resp, err := client.OrderNew(ctx, &req)
	if err != nil {
		t.Fatalf("create order fail %s", err.Error())
	}

	if resp.ID.String() != "9596912" || !resp.Amount.Equal(decimal.NewFromInt32(31431)) || !resp.Price.Equal(decimal.NewFromFloat(0.306525)) ||
		resp.Symbol.String() != "XRP-PERP" || resp.Status != exchange.OrderStatusOpen || resp.Type != exchange.OrderTypeLimit ||
		resp.Side != exchange.OrderSideSell {
		t.Errorf("bad order %v", *resp)
	}
}

func TestOrderFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	data := []byte(`
	{
  "success": true,
  "result": {
    "createdAt": "2019-03-05T09:56:55.728933+00:00",
    "filledSize": 10,
    "future": "XRP-PERP",
    "id": 9596912,
    "market": "XRP-PERP",
    "price": 0.306525,
    "avgFillPrice": 0.306526,
    "remainingSize": 31421,
    "side": "sell",
    "size": 31431,
    "status": "closed",
    "type": "limit",
    "reduceOnly": false,
    "ioc": false,
    "postOnly": false,
    "clientId": null
  }
}`)

	httpmock.RegisterResponder(http.MethodGet, "https://ftx.com/api/orders/9596912", httpmock.NewBytesResponder(200, data))

	ctx := context.Background()
	client := NewRestClient("", "")
	xrpS := newSwapSymbol("XRP")
	symbolMap["XRP-PERP"] = xrpS

	order, err := client.OrderFetch(ctx, &exchange.Order{
		ID: exchange.NewIntID(9596912),
	})
	if err != nil {
		t.Fatalf("cancel order fail %s", err.Error())
	}

	if order.Status != exchange.OrderStatusCancel || !order.AvgPrice.Equal(decimal.NewFromFloat(0.306526)) ||
		!order.Amount.Equal(decimal.NewFromInt(31431)) || !order.Filled.Equal(decimal.NewFromInt(10)) ||
		!order.Created.Equal(time.Date(2019, 3, 5, 9, 56, 55, 728933000, time.UTC)) {
		t.Errorf("bad fetch order %v", *order)
	}
}
