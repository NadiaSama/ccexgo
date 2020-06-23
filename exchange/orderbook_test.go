package exchange

import "testing"

func TestOrderBook(t *testing.T) {
	notify := &OrderBookNotify{
		Symbol: nil,
		Bids:   []OrderElem{},
		Asks:   []OrderElem{},
	}
	ods := NewOrderBookDS(notify)

	ods.Update(&OrderBookNotify{
		Symbol: nil,
		Bids:   []OrderElem{{1.0, 2.0}, {0.0, 1.0}, {0.5, 1.0}},
		Asks:   []OrderElem{{10.0, 1.0}, {3.0, 1.0}},
	})

	book := ods.Snapshot()
	if len(book.Bids) != 2 || book.Bids[0].Price != 1.0 || book.Bids[1].Price != 0.5 ||
		book.Asks[0].Price != 3.0 || book.Asks[1].Price != 10.0 {
		t.Errorf("bad snapshot %v", *book)
	}
}
