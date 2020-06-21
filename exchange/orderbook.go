package exchange

import (
	"time"

	"github.com/NadiaSama/ccexgo/misc/float"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
)

type (
	OrderElem struct {
		Price  float64
		Amount float64
	}

	//OrderBookNotify change of current orderbook
	//OrderElem.Amount == 0 means delete
	OrderBookNotify struct {
		Symbol string
		Bids   []OrderElem
		Asks   []OrderElem
	}

	OrderBookDS struct {
		symbol  string
		bids    *btree.Tree
		asks    *btree.Tree
		updated time.Time
	}

	OrderBook struct {
		Symbol  string
		Bids    []OrderElem
		Asks    []OrderElem
		Created time.Time
	}
)

func NewOrderBookDS(symbol string, bids []OrderElem, asks []OrderElem) *OrderBookDS {
	return &OrderBookDS{
		symbol:  symbol,
		bids:    newBook(bids),
		asks:    newBook(asks),
		updated: time.Now(),
	}
}

func newBook(data []OrderElem) *btree.Tree {
	l := len(data)
	if l < 3 {
		l = 3
	}
	tree := btree.NewWith(l, utils.Float64Comparator)
	for _, depth := range data {
		if float.Equal(depth.Price, 0.0) {
			continue
		}
		tree.Put(depth.Price, depth.Amount)
	}
	return tree
}

func (ds *OrderBookDS) Update(notify *OrderBookNotify) {
	updateTree := func(dest *btree.Tree, src []OrderElem) {
		for _, elem := range src {
			if float.Equal(elem.Price, 0.0) {
				continue
			}

			if float.Equal(elem.Amount, 0.0) {
				if _, ok := dest.Get(elem.Price); !ok {
					continue
				}
				dest.Remove(elem.Price)
			} else {
				dest.Put(elem.Price, elem.Amount)
			}
		}
	}

	updateTree(ds.bids, notify.Bids)
	updateTree(ds.asks, notify.Asks)
	ds.updated = time.Now()
}

func (ds *OrderBookDS) Snapshot() *OrderBook {
	ret := &OrderBook{
		Symbol:  ds.symbol,
		Bids:    make([]OrderElem, ds.bids.Size()),
		Asks:    make([]OrderElem, ds.asks.Size()),
		Created: ds.updated,
	}

	biter := ds.bids.Iterator()
	biter.End()
	i := 0
	for biter.Prev() {
		ret.Bids[i].Price = biter.Key().(float64)
		ret.Bids[i].Amount = biter.Value().(float64)
		i++
	}

	aiter := ds.asks.Iterator()
	aiter.Begin()
	i = 0
	for aiter.Next() {
		ret.Asks[i].Price = aiter.Key().(float64)
		ret.Asks[i].Amount = aiter.Value().(float64)
		i++
	}

	return ret
}
