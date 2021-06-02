package okex5

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/pkg/errors"
)

type (
	OrderSide     string
	PosSide       string
	OrdType       string
	TDMode        string
	OrderState    string
	OrderCategory string

	CreateOrderReq struct {
		InstID     string    `json:"instId"`
		TDMode     TDMode    `json:"tdMode"`
		Ccy        string    `json:"ccy,omitempty"`
		ClOrderID  string    `json:"clOrderId,omitempty"`
		Tag        string    `json:"tag,omitempty"`
		Side       OrderSide `json:"side"`
		PosSide    PosSide   `json:"posSide,omitempty"`
		OrdType    OrdType   `json:"ordType"`
		Sz         string    `json:"sz"`
		Px         string    `json:"px,omitempty"`
		ReduecOnly bool      `json:"reduceOnly"`
	}

	CreateOrderResp struct {
		OrderID   string `json:"ordId"`
		ClOrderID string `json:"clOrdId"`
		Tag       string `json:"tag"`
		SCode     string `json:"sCode"`
		SMsg      string `json:"sMsg"`
	}

	CancelOrderReq struct {
		InstID  string `json:"instId"`
		OrdId   string `json:"ordId"`
		ClOrdID string `json:"clOrdId"`
	}

	CancelOrderResp struct {
		OrdID   string `json:"ordId"`
		ClOrdID string `json:"clOrdId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	}

	FetchOrderReq struct {
		InstID  string `json:"instId"`
		OrdID   string `json:"ordId"`
		ClOrdID string `json:"clOrdId"`
	}

	Order struct {
		InstType    string        `json:"instType"`
		InstId      string        `json:"instId"`
		Ccy         string        `json:"ccy"`
		OrderID     string        `json:"orderId"`
		ClOrdID     string        `json:"clOrdId"`
		Tag         string        `json:"tag"`
		Px          string        `json:"px"`
		Sz          string        `json:"sz"`
		Pnl         string        `json:"pnl"`
		OrderType   OrdType       `json:"ordType"`
		Side        OrderSide     `json:"side"`
		PosSide     PosSide       `json:"posSide"`
		TDMode      TDMode        `json:"tdMode"`
		AccFillSZ   string        `json:"accFillSz"`
		FillPx      string        `json:"fillPx"`
		TradeID     string        `json:"tradeId"`
		FillSz      string        `json:"fillSz"`
		FillTime    string        `json:"fillTime"`
		AvgPx       string        `json:"avgPx"`
		State       OrderState    `json:"state"`
		Lever       string        `json:"lever"`
		TpTriggerPx string        `json:"tpTriggerPx"`
		TpOrdPx     string        `json:"tpOrdPx"`
		SlTriggerPx string        `json:"slTriggerPx"`
		SlOrdPx     string        `json:"slOrdPx"`
		FeeCcy      string        `json:"feeCcy"`
		Fee         string        `json:"fee"`
		RebateCcy   string        `json:"rebatCcy"`
		Rebat       string        `json:"rebat"`
		Category    OrderCategory `json:"category"`
		UTime       string        `json:"uTime"`
		CTime       string        `json:"cTime"`
	}
)

const (
	TDModeIsolated TDMode = "isolated"
	TDModeCross    TDMode = "cross"
	TDModeCash     TDMode = "cash"

	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"

	PosSideNone  PosSide = ""
	PosSideLong  PosSide = "long"
	PosSideShort PosSide = "short"
	PosSideNet   PosSide = "net"

	OrdTypeMaket           OrdType = "market"
	OrdTypeLimit           OrdType = "limit"
	OrdTypePostOnly        OrdType = "post_only"
	OrdTypeFOK             OrdType = "fok"
	OrdTypeIOC             OrdType = "ioc"
	OrdTypeOptimalLimitIOC OrdType = "optimal_limit_ioc"

	OrderStateCanceled        OrderState = "canceled"
	OrderStateLive            OrderState = "live"
	OrderStatePartiallyFilled OrderState = "partially_filled"
	OrderStateFilled          OrderState = "filled"

	OrderCategoryUnknown            OrderCategory = ""
	OrderCategoryNormal             OrderCategory = "normal"
	OrderCategoryTwap               OrderCategory = "twap"
	OrderCategoryAdl                OrderCategory = "adl"
	OrderCategoryFullLiquidation    OrderCategory = "full_liquidation"
	OrderCategoryPartialLiquidation OrderCategory = "partial_liquidation"
	OrderCategoryDelivery           OrderCategory = "delivery"
)

const (
	CreateOrderEndPoint = "/api/v5/trade/order"
	FetchOrderEndPoint  = CreateOrderEndPoint
	CancelOrderEndPoint = "/api/v5/trade/cancel-order"
)

var (
	orderTypeMap     map[string]OrdType       = make(map[string]OrdType)
	orderSideMap     map[string]OrderSide     = make(map[string]OrderSide)
	posSideMap       map[string]PosSide       = make(map[string]PosSide)
	tdModeMap        map[string]TDMode        = make(map[string]TDMode)
	orderCategoryMap map[string]OrderCategory = make(map[string]OrderCategory)
	orderStateMap    map[string]OrderState    = make(map[string]OrderState)
)

func init() {
	ots := []OrdType{
		OrdTypeMaket,
		OrdTypeLimit,
		OrdTypePostOnly,
		OrdTypeFOK,
		OrdTypeIOC,
		OrdTypeOptimalLimitIOC,
	}

	for _, t := range ots {
		orderTypeMap[string(t)] = t
	}

	ors := []OrderSide{
		OrderSideBuy,
		OrderSideSell,
	}
	for _, t := range ors {
		orderSideMap[string(t)] = t
	}

	pss := []PosSide{
		PosSideLong,
		PosSideNone,
		PosSideShort,
		PosSideNet,
	}
	for _, p := range pss {
		posSideMap[string(p)] = p
	}

	tss := []TDMode{
		TDModeCash,
		TDModeCross,
		TDModeIsolated,
	}
	for _, t := range tss {
		tdModeMap[string(t)] = t
	}

	cs := []OrderCategory{
		OrderCategoryUnknown,
		OrderCategoryAdl,
		OrderCategoryNormal,
		OrderCategoryDelivery,
		OrderCategoryFullLiquidation,
		OrderCategoryPartialLiquidation,
		OrderCategoryTwap,
	}
	for _, t := range cs {
		orderCategoryMap[string(t)] = t
	}

	oss := []OrderState{
		OrderStateCanceled,
		OrderStateFilled,
		OrderStateLive,
		OrderStatePartiallyFilled,
	}
	for _, t := range oss {
		orderStateMap[string(t)] = t
	}
}

func (rc *RestClient) CreateOrder(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	ret := []CreateOrderResp{}
	if err := rc.doPostJSON(ctx, CreateOrderEndPoint, req, &ret); err != nil {
		return nil, err
	}

	return &ret[0], nil
}

func (rc *RestClient) CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, error) {
	ret := []CancelOrderResp{}
	if err := rc.doPostJSON(ctx, CancelOrderEndPoint, req, &ret); err != nil {
		return nil, err
	}

	return &ret[0], nil
}

func (rc *RestClient) FetchOrder(ctx context.Context, req *FetchOrderReq) (*Order, error) {
	ret := []Order{}
	values := url.Values{}
	values.Add("instId", req.InstID)
	if len(req.OrdID) != 0 {
		values.Add("ordId", req.OrdID)
	} else if len(req.ClOrdID) != 0 {
		values.Add("clOrdId", req.ClOrdID)
	} else {
		return nil, errors.Errorf("ordID or clOrdId is required")
	}
	if err := rc.Request(ctx, http.MethodGet, FetchOrderEndPoint, values, nil, true, &ret); err != nil {
		return nil, err
	}

	return &ret[0], nil
}

func (rc *RestClient) doPostJSON(ctx context.Context, endPoint string, obj interface{}, dst interface{}) error {
	raw, err := json.Marshal(obj)
	if err != nil {
		return errors.WithMessage(err, "marshal json error")
	}

	body := bytes.NewBuffer(raw)
	if err := rc.Request(ctx, http.MethodPost, endPoint, nil, body, true, dst); err != nil {
		return err
	}
	return nil
}

func (ot *OrdType) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(orderTypeMap, "orderType", raw, ot)
}

func (os *OrderSide) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(orderSideMap, "orderSide", raw, os)
}

func (tm *TDMode) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(tdModeMap, "tdMode", raw, tm)
}

func (os *OrderState) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(orderStateMap, "orderState", raw, os)
}

func (oc *OrderCategory) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(orderCategoryMap, "orderCategory", raw, oc)
}

func (ps *PosSide) UnmarshalJSON(raw []byte) error {
	var key string
	if err := json.Unmarshal(raw, &key); err != nil {
		return errors.Errorf("invalid key %s", string(raw))
	}

	val, ok := posSideMap[key]
	if !ok {
		return errors.Errorf("unkown posSide %s", key)
	}
	*ps = val
	return nil
}

func assignMapPtr(dict interface{}, typName string, rawKey []byte, dst interface{}) error {
	var key string
	if err := json.Unmarshal(rawKey, &key); err != nil {
		return errors.Errorf("invalid key %s", string(rawKey))
	}

	kVal := reflect.ValueOf(key)
	dictVal := reflect.ValueOf(dict)

	sVal := dictVal.MapIndex(kVal)
	if !sVal.IsValid() {
		return errors.Errorf("unkown %s '%s'", typName, string(key))
	}

	dVal := reflect.Indirect(reflect.ValueOf(dst))
	dVal.Set(sVal)
	return nil
}
