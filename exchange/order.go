package exchange

import "time"

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
		Price  float64
		Amount float64
	}

	Order struct {
		ID       OrderID
		Symbol   Symbol
		Amount   float64
		Filled   float64
		Price    float64
		AvgPrice float64
		Fee      float64
		Created  time.Time
		Updated  time.Time
		Side     OrderSide
		Status   OrderStatus
		Type     OrderType
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
