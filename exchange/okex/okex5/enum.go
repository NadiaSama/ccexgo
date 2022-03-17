package okex5

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	CodeOK = "0"
)

type (
	OrderSide     string
	PosSide       string
	OrdType       string
	TDMode        string
	OrderState    string
	OrderCategory string
	InstType      string
	ExecType      string
	MgnMode       string
	CtType        string
)

const (
	CtTypeNone    CtType = ""
	CtTypeLinear  CtType = "linear"
	CtTypeInverse CtType = "Inverse"

	MgnModeNone     MgnMode = ""
	MgnModeCash     MgnMode = "cash"
	MgnModeCross    MgnMode = "cross"
	MgnModeIsolated MgnMode = "isolated"
	MgnModeEmpty    MgnMode = ""

	ExecTypeMaker ExecType = "M"
	ExecTypeTaker ExecType = "T"

	InstTypeSpot    InstType = "SPOT"
	InstTypeMargin  InstType = "MARGIN"
	InstTypeSwap    InstType = "SWAP"
	InstTypeFutures InstType = "FUTURES"
	InstTypeOption  InstType = "OPTION"
	InstTypeAny     InstType = "ANY"
	InstTypeNone    InstType = ""

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

var (
	orderTypeMap     map[string]OrdType       = make(map[string]OrdType)
	orderSideMap     map[string]OrderSide     = make(map[string]OrderSide)
	posSideMap       map[string]PosSide       = make(map[string]PosSide)
	tdModeMap        map[string]TDMode        = make(map[string]TDMode)
	orderCategoryMap map[string]OrderCategory = make(map[string]OrderCategory)
	orderStateMap    map[string]OrderState    = make(map[string]OrderState)
	instTypeMap      map[string]InstType      = make(map[string]InstType)
	execTypeMap      map[string]ExecType      = make(map[string]ExecType)
	mgnModeMap       map[string]MgnMode       = make(map[string]MgnMode)
	ctTypeMap        map[string]CtType        = make(map[string]CtType)
)

func (it *InstType) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(instTypeMap, "instType", raw, it)
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
	return assignMapPtr(posSideMap, "posSide", raw, ps)
}

func (et *ExecType) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(execTypeMap, "execType", raw, et)
}

func (mgn *MgnMode) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(mgnModeMap, "mgnMode", raw, mgn)
}

func (ct *CtType) UnmarshalJSON(raw []byte) error {
	return assignMapPtr(ctTypeMap, "ctType", raw, ct)
}

func init() {
	its := []InstType{
		InstTypeSpot,
		InstTypeMargin,
		InstTypeFutures,
		InstTypeSwap,
		InstTypeOption,
		InstTypeAny,
		InstTypeNone,
	}
	for _, i := range its {
		instTypeMap[string(i)] = i
	}

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

	ess := []ExecType{
		ExecTypeMaker,
		ExecTypeTaker,
	}
	for _, e := range ess {
		execTypeMap[string(e)] = e
	}

	cts := []CtType{
		CtTypeNone,
		CtTypeInverse,
		CtTypeLinear,
	}
	for _, c := range cts {
		ctTypeMap[string(c)] = c
	}

	mms := []MgnMode{
		MgnModeNone,
		MgnModeCash,
		MgnModeCross,
		MgnModeIsolated,
		MgnModeEmpty,
	}
	for _, m := range mms {
		mgnModeMap[string(m)] = m
	}
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
