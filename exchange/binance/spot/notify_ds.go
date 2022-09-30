package spot

import (
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/tidwall/gjson"
)

type (

	// BookTickerNotify binance spot bookticker notify
	BookTickerNotify struct {
		UpdateID   int64  `json:"u"`
		Symbol     string `json:"s"`
		Bid1Price  string `json:"b"`
		Bid1Amount string `json:"B"`
		Ask1Price  string `json:"a"`
		Ask1Amount string `json:"A"`
	}

	// BookTickerChannel xxx@bookTicker channel
	BookTickerChannel struct {
		symbol string
	}
)

func NewTickerChannel(symbol string) exchange.Channel {
	return &BookTickerChannel{
		symbol: strings.ToLower(symbol),
	}
}

func (tc *BookTickerChannel) String() string {
	return fmt.Sprintf("%s@bookTicker", tc.symbol)
}

func ParseBookTickerNotify(g *gjson.Result) *BookTickerNotify {

	updateID := g.Get("u").Int()
	symbol := g.Get("s").String()
	bid1Price := g.Get("b").String()
	bid1Amount := g.Get("B").String()
	ask1Price := g.Get("a").String()
	ask1Amount := g.Get("A").String()

	tn := &BookTickerNotify{
		UpdateID:   updateID,
		Symbol:     symbol,
		Bid1Price:  bid1Price,
		Bid1Amount: bid1Amount,
		Ask1Price:  ask1Price,
		Ask1Amount: ask1Amount,
	}
	return tn
}
