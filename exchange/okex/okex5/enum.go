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
