package swap

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	orderTable = "swap/order"
)

type (
	OrderChannel struct {
		sym exchange.SwapSymbol
	}

	Order struct {
		InstrumentID string          `json:"instrument_id"`
		Size         decimal.Decimal `json:"size"`
		Timestamp    string          `json:"timestamp"`
		FilledQty    decimal.Decimal `json:"filled_qty"`
		Fee          decimal.Decimal `json:"fee"`
		OrderID      string          `json:"order_id"`
		ClientID     string          `json:"client_id"`
		Price        decimal.Decimal `json:"price"`
		PriceAvg     decimal.Decimal `json:"price_avg"`
		Type         string          `json:"type"`
		ContractVal  decimal.Decimal `json:"contract_val"`
		OrderType    string          `json:"order_type"`
		State        string          `json:"state"`
	}
)

var (
	sideMap map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"1": exchange.OrderSideBuy,
		"2": exchange.OrderSideSell,
		"3": exchange.OrderSideCloseLong,
		"4": exchange.OrderSideCloseShort,
	}

	statusMap map[string]exchange.OrderStatus = map[string]exchange.OrderStatus{
		"-1": exchange.OrderStatusCancel,
		"0":  exchange.OrderStatusOpen,
		"1":  exchange.OrderStatusOpen,
		"2":  exchange.OrderStatusDone,
		"3":  exchange.OrderStatusOpen,
		"4":  exchange.OrderStatusOpen,
	}
)

func NewOrderChannel(symbol exchange.SwapSymbol) exchange.Channel {
	return &OrderChannel{
		sym: symbol,
	}
}

func (oc *OrderChannel) String() string {
	return fmt.Sprintf("%s:%s", orderTable, oc.String())
}

func init() {
	okex.SubscribeCB(orderTable, nil)
}

func parseOrderCB(table string, action string, raw json.RawMessage) (*rpc.Notify, error) {
	var o Order
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}

	order, err := o.Transform()
	if err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: orderTable,
		Params: order,
	}, nil
}

func (o *Order) Transform() (*exchange.Order, error) {
	sym, err := okex.ParseSwapSymbol(o.InstrumentID)
	if err != nil {
		return nil, err
	}

	side, ok := sideMap[o.Type]
	if !ok {
		return nil, errors.Errorf("unkown order side '%s'", o.Type)
	}

	status, ok := statusMap[o.State]
	if !ok {
		return nil, errors.Errorf("unkown order state '%s'", o.State)
	}

	ts, err := okex.ParseTime(o.Timestamp)
	if err != nil {
		return nil, err
	}

	return &exchange.Order{
		Symbol:   sym,
		Side:     side,
		ID:       exchange.NewStrID(o.OrderID),
		ClientID: exchange.NewStrID(o.ClientID),
		Status:   status,
		Amount:   o.Size,
		Filled:   o.FilledQty,
		Price:    o.Price,
		AvgPrice: o.PriceAvg,
		Fee:      o.Fee,
		Updated:  ts,
		Raw:      o,
	}, nil
}
