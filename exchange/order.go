package exchange

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type (
	OrderSide   int
	OrderType   int
	OrderStatus int

	OrderID interface {
		String() string
	}

	//OrderRequest carry field which used to create order
	OrderRequest struct {
		Symbol Symbol
		Side   OrderSide
		Type   OrderType
		Price  decimal.Decimal
		Amount decimal.Decimal
	}

	Order struct {
		ID       OrderID
		Symbol   Symbol
		Amount   decimal.Decimal
		Filled   decimal.Decimal
		Price    decimal.Decimal
		AvgPrice decimal.Decimal
		Fee      decimal.Decimal
		Created  time.Time
		Updated  time.Time
		Side     OrderSide
		Status   OrderStatus
		Type     OrderType
		Raw      interface{} `json:"omitempty"`
	}

	//OrderReqOption specific option to create order
	//each exchange support different options.
	OrderReqOption interface {
	}

	//PostOnlyOption wether the order ensure maker
	PostOnlyOption struct {
		PostOnly bool
	}

	//TimeInForceFlag specific TimeInForceOption value
	TimeInForceFlag string
	//TimeInForceOption specific how long the order
	//remains in effect
	TimeInForceOption struct {
		Flag TimeInForceFlag
	}

	IntID struct {
		ID int64
	}
)

const (
	OrderSideBuy = iota
	OrderSideSell

	OrderTypeLimit = iota
	OrderTypeMarket
	OrderTypeStopLimit
	OrderTypeStopMarket

	//OrderStatusUnknown means order info need check with api
	OrderStatusUnknown = iota
	OrderStatusOpen
	OrderStatusDone
	OrderStatusCancel
	OrderStatusFailed

	//TimeInForceGTC good_till_cancel
	TimeInForceGTC = "gtc"
	//TimeInForceFOK fill_or_kill
	TimeInForceFOK = "fok"
	//TimeInForceIOC immediate_or_cancel
	TimeInForceIOC = "ioc"
)

func NewPostOnlyOption(postOnly bool) OrderReqOption {
	return &PostOnlyOption{
		PostOnly: postOnly,
	}
}

func NewTimeInForceOption(flag TimeInForceFlag) OrderReqOption {
	return &TimeInForceOption{
		Flag: flag,
	}
}

func NewIntID(id int64) IntID {
	return IntID{id}
}

func (sid IntID) String() string {
	return strconv.FormatInt(sid.ID, 10)
}
