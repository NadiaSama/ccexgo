package spot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	TradeDetailChannel struct {
		sym string
	}

	TradeDetailData struct {
		ID        int64   `json:"id"`
		TS        int64   `json:"ts"`
		TradeID   int64   `json:"tradeId"`
		Amount    float64 `json:"amount"`
		Price     float64 `json:"price"`
		Direction string  `json:"direction"`
	}

	TradeDetail struct {
		ID   int64             `json:"id"`
		TS   int64             `json:"ts"`
		Data []TradeDetailData `json:"data"`
	}
)

var (
	HuobiSide2ExchangeSide map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"buy":  exchange.OrderSideBuy,
		"sell": exchange.OrderSideSell,
	}
)

func IsTradeDetailChanel(ch string) bool {
	ss := strings.Split(ch, ".")
	if len(ss) != 4 {
		return false
	}

	return ss[0] == "market" && ss[2] == "trade" && ss[3] == "detail"
}

func TradeDetailSymbol(ch string) string {
	ss := strings.Split(ch, ".")
	return ss[1]
}

func NewTradeDetailChanel(sym string) *TradeDetailChannel {
	return &TradeDetailChannel{
		sym: strings.ToLower(sym),
	}
}

func (tdc *TradeDetailChannel) String() string {
	return fmt.Sprintf("maket.%s.trade.detail", tdc.sym)
}

func ParseTradeTick(ch string, ts int64, raw json.RawMessage) ([]*exchange.Trade, error) {
	sym := TradeDetailSymbol(ch)
	symbol, err := ParseSymbol(sym)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	var td TradeDetail
	if err := json.Unmarshal(raw, &td); err != nil {
		return nil, errors.WithMessage(err, "parse trade detail fail")
	}

	return td.Parse(symbol)
}

func (td *TradeDetail) Parse(sym exchange.Symbol) ([]*exchange.Trade, error) {
	var ret []*exchange.Trade

	for i, ch := range td.Data {
		side, ok := HuobiSide2ExchangeSide[ch.Direction]
		if !ok {
			return nil, errors.Errorf("unsupport direction=%s", ch.Direction)
		}

		trade := &exchange.Trade{
			ID:     strconv.FormatInt(ch.ID, 10),
			Time:   huobi.ParseTS(ch.TS),
			Symbol: sym,
			Price:  decimal.NewFromFloat(ch.Price),
			Amount: decimal.NewFromFloat(ch.Amount),
			Side:   side,
			Raw:    &td.Data[i],
		}

		ret = append(ret, trade)
	}
	return ret, nil
}
