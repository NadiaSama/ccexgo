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
		Symbol   Symbol
		ClientID OrderID
		Side     OrderSide
		Type     OrderType
		Price    decimal.Decimal
		Amount   decimal.Decimal
	}

	Order struct {
		ID       OrderID
		ClientID OrderID
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
		Raw      interface{} `json:"-"`
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

	StrID string
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

func NewOrderRequest(sym Symbol, cid OrderID, side OrderSide, typ OrderType,
	price float64, amount float64) *OrderRequest {
	return &OrderRequest{
		Symbol:   sym,
		ClientID: cid,
		Side:     side,
		Type:     typ,
		Price:    decimal.NewFromFloat(price),
		Amount:   decimal.NewFromFloat(amount),
	}
}

func NewStrID(id string) StrID {
	return StrID(id)
}
func (sid StrID) String() string {
	return string(sid)
}

func NewIntID(id int64) IntID {
	return IntID{id}
}
func (sid IntID) String() string {
	return strconv.FormatInt(sid.ID, 10)
}

//Equal check whether o equal o2 mainly used for test
func (o *Order) Equal(o2 *Order) bool {
	return o.ID.String() == o2.ID.String() &&
		o.Amount.Equal(o2.Amount) && o.Filled.Equal(o2.Filled) &&
		o.Price.Equal(o2.Price) && o.AvgPrice.Equal(o2.AvgPrice) && o.Fee.Equal(o2.Fee) &&
		o.Created.Equal(o2.Created) && o.Updated.Equal(o2.Updated) &&
		o.Status == o2.Status && o.Type == o2.Type && o.Side == o2.Side &&
		o.Symbol.String() == o2.Symbol.String()
}
