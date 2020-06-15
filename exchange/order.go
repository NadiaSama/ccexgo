package exchange

import "time"

type (
	OrderSide   int
	OrderType   int
	OrderStatus int

	Order struct {
		ID       interface{}
		Symbol   Symbol
		Amount   float64
		Price    float64
		AvgPrice float64
		Fee      float64
		Created  time.Time
		Updated  time.Time
		Side     OrderSide
		Status   OrderStatus
	}
)

const (
	OrderSideBuy = iota
	OrderSideSell

	OrderTypeLimit = iota
	OrderTypeMarket
	OrderTypeStopLimit
	OrderTypeStopMarket

	OrderStatusOpen = iota
	OrderStatusDone
	OrderStatusFilled
	OrderStatusCancel
)
