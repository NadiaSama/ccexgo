package swap

import (
	"context"
	"encoding/json"
	"strings"
)

type (
	OrderReq struct {
		data map[string]interface{}
	}

	OrderResp struct {
		OrderID    int64  `json:"order_id"`
		OrderIDStr string `json:"order_id_str"`
	}

	SwapCancelReq struct {
		symbol         string
		orderIDS       []string
		clientOrderIDS []string
	}

	SwapCancelError struct {
		OrderID string `json:"order_id"`
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
	}

	SwapCancelResp struct {
		Errors    []SwapCancelError `json:"errors"`
		Successes string            `json:"successes"` //id1,id2,id3 ...
	}
)

const (
	SwapOrderEndPoint  = "/swap-api/v1/swap_order"
	SwapCancelEndPoint = "/swap-api/v1/swap_cancel"

	OrderDirectionBuy  = "buy"
	OrderDirectionSell = "sell"
	OrderOffsetOpen    = "open"
	OrderOffsetClose   = "close"

	OrderPriceLimit  = "limit"
	OrderPriceMarket = "opponent"
)

func NewOrderReq(contractCode string, volume int, direction string, offset string, lever int, orderPriceType string) *OrderReq {
	ret := OrderReq{
		data: make(map[string]interface{}),
	}

	ret.data["contract_code"] = contractCode
	ret.data["volume"] = volume
	ret.data["direction"] = direction
	ret.data["offset"] = offset
	ret.data["lever_rate"] = lever
	ret.data["order_price_type"] = orderPriceType
	return &ret
}

func (or *OrderReq) Price(price float64) *OrderReq {
	or.data["price"] = price
	return or
}

func (or *OrderReq) Serialize() ([]byte, error) {
	return json.Marshal(or.data)
}

func NewSwapCancelReq(symbol string) *SwapCancelReq {
	return &SwapCancelReq{
		symbol: symbol,
	}
}

func (scr *SwapCancelReq) Orders(ids ...string) *SwapCancelReq {
	for _, id := range ids {
		scr.orderIDS = append(scr.orderIDS, id)
	}
	return scr
}

func (scr *SwapCancelReq) ClientOrderIDs(ids ...string) *SwapCancelReq {
	for _, id := range ids {
		scr.clientOrderIDS = append(scr.clientOrderIDS, id)
	}
	return scr
}

func (scr *SwapCancelReq) Serialize() ([]byte, error) {
	data := map[string]string{
		"contract_code": scr.symbol,
	}

	if len(scr.orderIDS) != 0 {
		data["order_id"] = strings.Join(scr.orderIDS, ",")
	}

	if len(scr.clientOrderIDS) != 0 {
		data["client_order_id"] = strings.Join(scr.clientOrderIDS, ",")
	}

	return json.Marshal(data)
}

func (rc *RestClient) SwapOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var resp OrderResp
	if err := rc.PrivatePostReq(ctx, SwapOrderEndPoint, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (rc *RestClient) SwapCancel(ctx context.Context, req *SwapCancelReq) (*SwapCancelResp, error) {
	var resp SwapCancelResp
	if err := rc.PrivatePostReq(ctx, SwapCancelEndPoint, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
