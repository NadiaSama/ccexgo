package swap

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	DepthHighFreqChannel struct {
		contractCode string
		size         int
	}

	Depth struct {
		Asks    [][2]float64
		Bids    [][2]float64
		Ch      string
		Event   string
		ID      int64
		MrID    int64
		TS      int64
		Version int
	}
)

func NewDepthHighFreq(symbol exchange.SwapSymbol, size int) exchange.Channel {
	return &DepthHighFreqChannel{
		contractCode: symbol.String(),
		size:         size,
	}
}

func (ch *DepthHighFreqChannel) String() string {
	return fmt.Sprintf("market.%s.depth.size_%d.high_freq", ch.contractCode, ch.size)
}

func ParseDepth(raw json.RawMessage) (interface{}, error) {
	var d Depth
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}

	fields := strings.Split(d.Ch, ".")
	if len(fields) < 2 {
		return nil, errors.Errorf("invalid ch %s", d.Ch)
	}

	sym, err := ParseSymbol(fields[1])
	if err != nil {
		return nil, err
	}

	bids := make([]exchange.OrderElem, len(d.Bids))
	asks := make([]exchange.OrderElem, len(d.Asks))

	for i, b := range d.Bids {
		bids[i] = exchange.OrderElem{
			Price:  b[0],
			Amount: b[1],
		}
	}

	for i, a := range d.Asks {
		asks[i] = exchange.OrderElem{
			Price:  a[0],
			Amount: a[1],
		}
	}

	return &exchange.OrderBook{
		Symbol:  sym,
		Bids:    bids,
		Asks:    asks,
		Created: time.Unix(d.TS/1000, d.TS%1000*1e6),
		Raw:     &d,
	}, nil
}
