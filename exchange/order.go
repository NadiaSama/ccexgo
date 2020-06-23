package exchange

import "time"

type (
	OrderSide   int
	OrderType   int
	OrderStatus int
	OrderID     interface{}

	//OrderRequest carry field which used to create order
	OrderRequest struct {
		Symbol Symbol
		Side   OrderSide
		Type   OrderType
		Price  float64
		Amount float64
		Opt    interface{}
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
)
