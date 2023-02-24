package swap

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type (
	BookTickerChannel struct {
		symbol string
	}

	BookTickerNotify struct {
		Event      string `json:"e"`
		UpdateID   int64  `json:"u"`
		PushTime   int64  `json:"E"`
		MatchTime  int64  `json:"T"`
		Symbol     string `json:"s"`
		Bid1Price  string `json:"b"`
		Bid1Amount string `json:"B"`
		Ask1Price  string `json:"a"`
		Ask1Amount string `json:"A"`
	}
)

func NewBookTickerChannel(sym string) *BookTickerChannel {
	return &BookTickerChannel{
		symbol: strings.ToLower(sym),
	}
}

func (btc *BookTickerChannel) String() string {
	return fmt.Sprintf("%s@bookTicker", btc.symbol)
}

func ParseBookTickerNotify(g *gjson.Result) *BookTickerNotify {
	event := g.Get("e").String()
	UpdateID := g.Get("u").Int()
	pushTime := g.Get("E").Int()
	matchTime := g.Get("T").Int()
	symbol := g.Get("s").String()
	bid1Price := g.Get("b").String()
	bid1Amount := g.Get("B").String()
	ask1Price := g.Get("a").String()
	ask1Amount := g.Get("A").String()

	return &BookTickerNotify{
		Event:      event,
		UpdateID:   UpdateID,
		PushTime:   pushTime,
		MatchTime:  matchTime,
		Symbol:     symbol,
		Bid1Price:  bid1Price,
		Bid1Amount: bid1Amount,
		Ask1Price:  ask1Price,
		Ask1Amount: ask1Amount,
	}
}
