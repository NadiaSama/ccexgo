package spot

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"
)

type (
	MBPChannel struct {
		symbol exchange.Symbol
		size   int
	}

	Depth struct {
		SeqNum     int64        `json:"seqNum"`
		PrevSeqNum int64        `json:"prevSeqNum"`
		Bids       [][2]float64 `json:"bids"`
		Asks       [][2]float64 `json:"asks"`
	}

	MBPFullReq struct {
		Req string `json:"req"`
		ID  string `json:"id"`
	}

	MBPFullResp struct {
		ID     string
		Status string
		TS     int64
		Req    string
		Data   Depth
	}

	//MBPDepthDS build depth according incremental updates and refresh message
	MBPDepthDS struct {
		bids       *btree.Tree
		asks       *btree.Tree
		cache      []Depth
		refresh    *Depth
		lastSeqNum int64
		ts         time.Time
		inited     bool
		symbol     exchange.Symbol
	}
)

func NewMBPFullReq(symbol exchange.Symbol, size int) *MBPFullReq {
	return &MBPFullReq{
		Req: fmt.Sprintf("market.%s.mbp.%d", symbol.String(), size),
		ID:  "123",
	}
}

func NewMBPChannel(sym exchange.Symbol, size int) exchange.Channel {
	return &MBPChannel{
		symbol: sym,
		size:   size,
	}
}

func (mc *MBPChannel) String() string {
	return fmt.Sprintf("market.%s.mbp.%d", mc.symbol.String(), mc.size)
}

func NewMBPDepthDS(symbol exchange.Symbol) *MBPDepthDS {
	return &MBPDepthDS{
		bids:   btree.NewWith(3, utils.Float64Comparator),
		asks:   btree.NewWith(3, utils.Float64Comparator),
		cache:  make([]Depth, 16),
		symbol: symbol,
	}
}

//Push add incremental updates into ds
func (ds *MBPDepthDS) Push(d *Depth, ts time.Time) (inited bool, err error) {
	//make sure seqNum is consistent
	if ds.lastSeqNum != 0 && ds.lastSeqNum != d.PrevSeqNum {
		err = errors.Errorf("seqNum inconsistend lastSeqNum=%d PrevSeqNum=%d", ds.lastSeqNum, d.PrevSeqNum)
		return
	}
	ds.lastSeqNum = d.SeqNum
	ds.ts = ts

	if !ds.inited && ds.refresh == nil {
		ds.cache = append(ds.cache, *d)
		inited = false
		return
	}

	if ds.inited {
		processArr(ds.bids, d.Bids)
		processArr(ds.asks, d.Asks)
		inited = true
		return
	}

	ds.cache = append(ds.cache, *d)
	if d.PrevSeqNum >= ds.refresh.SeqNum {
		ds.init()
	}
	return ds.inited, nil
}

func (ds *MBPDepthDS) AddRefresh(d *Depth) {
	ds.refresh = d

	if ds.cache[len(ds.cache)-1].PrevSeqNum < ds.refresh.SeqNum {
		return
	}

	ds.init()
}

//OrderBook generate orderbook according ds bids, asks structure
//size specific orderbook size.  -1 means use bids, asks size
func (ds *MBPDepthDS) OrderBook(size int) *exchange.OrderBook {
	var bidLen, askLen int
	if size > ds.bids.Size() || size == -1 {
		bidLen = ds.bids.Size()
	} else {
		bidLen = size
	}

	if size > ds.asks.Size() || size == -1 {
		askLen = ds.asks.Size()
	} else {
		askLen = size
	}

	ret := &exchange.OrderBook{
		Symbol:  ds.symbol,
		Created: ds.ts,
		Bids:    make([]exchange.OrderElem, bidLen),
		Asks:    make([]exchange.OrderElem, askLen),
	}

	iter := ds.bids.Iterator()
	iter.End()
	for i := 0; i < bidLen; i++ {
		iter.Prev()
		ret.Bids[i] = exchange.OrderElem{
			Price:  iter.Key().(float64),
			Amount: iter.Value().(float64),
		}
	}

	iter = ds.asks.Iterator()
	for i := 0; i < askLen; i++ {
		iter.Next()
		ret.Asks[i] = exchange.OrderElem{
			Price:  iter.Key().(float64),
			Amount: iter.Value().(float64),
		}
	}

	return ret
}

func processArr(tr *btree.Tree, elems [][2]float64) {
	for _, e := range elems {
		if e[1] == 0.0 {
			tr.Remove(e[0])
			continue
		}

		tr.Put(e[0], e[1])
	}
}

func (ds *MBPDepthDS) init() {
	processArr(ds.bids, ds.refresh.Bids)
	processArr(ds.asks, ds.refresh.Asks)

	idx := sort.Search(len(ds.cache), func(i int) bool {
		return ds.cache[i].PrevSeqNum >= ds.refresh.SeqNum
	})

	for i := idx; i < len(ds.cache); i++ {
		d := ds.cache[i]
		processArr(ds.bids, d.Bids)
		processArr(ds.asks, d.Asks)
	}

	ds.inited = true
}

func ParseDepth(ch string, ts int64, tick json.RawMessage) (interface{}, error) {
	var d Depth
	if err := json.Unmarshal(tick, &d); err != nil {
		return nil, err
	}

	fields := strings.Split(ch, ".")
	if len(fields) < 2 {
		return nil, errors.Errorf("invalid channedl %s", ch)
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
		Created: time.Unix(ts/1e3, ts%1e3*1e5),
		Bids:    bids,
		Asks:    asks,
		Raw:     &d,
	}, nil
}
