package okex5

import (
	"encoding/json"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
)

type (
	//RawDepth incremental depth push
	RawDepth struct {
		Asks     [][4]string `json:"asks"`
		Bids     [][4]string `json:"bids"`
		Ts       string      `json:"ts"`
		Checksum int32       `json:"checksum"`
	}

	//DepthDS recv okex5 RawDepth notify and calc depth
	DepthDS struct {
		bids    *btree.Tree
		asks    *btree.Tree
		updated time.Time
	}

	Depth struct {
		Asks         [][4]string
		Bids         [][4]string
		Ts           string
		Checksum     int32 //checksum recv from okex
		CalcChecskum int32 //checksum calc by bids and asks data
	}
)

const (
	DepthSnapshot = "snapshot"
	DepthUpdate   = "update"

	Books5Channel     = "books5"
	Books50TBTChannel = "books50-l2-tbt"
)

func init() {
	chs := []string{Books5Channel, Books50TBTChannel}

	for _, c := range chs {
		parseCBMap[c] = parseDepth
	}
}

func NewBooks5Channel(instId string) exchange.Channel {
	return &Okex5Channel{
		InstID:  instId,
		Channel: Books5Channel,
	}
}

func NewBooks50TBTChannel(instId string) exchange.Channel {
	return &Okex5Channel{
		InstID:  instId,
		Channel: Books50TBTChannel,
	}
}

func NewDepthDS() *DepthDS {
	ret := &DepthDS{
		bids: btree.NewWith(32, floatComparator),
		asks: btree.NewWith(32, floatComparator),
	}

	return ret
}

func parseDepth(data *wsResp) (*rpc.Notify, error) {
	var d []RawDepth

	if err := json.Unmarshal(data.Data, &d); err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: data.Arg.Channel,
		Params: d,
	}, nil
}

//Push update depth data accoring raw depth data base on the
//https://www.okex.com/docs-v5/en/#websocket-api-checksum-merging-incremental-data-into-full-data
func (ds *DepthDS) Push(raw *RawDepth) (*Depth, error) {
	updateBook(ds.asks, raw.Asks)
	updateBook(ds.bids, raw.Bids)

	ts, err := ParseTimestamp(raw.Ts)
	if err != nil {
		return nil, err
	}
	ds.updated = ts

	ret := ds.snapShot()
	ret.Ts = raw.Ts
	ret.Checksum = raw.Checksum

	return ret, nil
}

func (ds *DepthDS) snapShot() *Depth {
	ret := &Depth{
		Bids: make([][4]string, ds.bids.Size()),
		Asks: make([][4]string, ds.asks.Size()),
	}

	biter := ds.bids.Iterator()
	biter.End()
	i := 0
	for biter.Prev() {
		ret.Bids[i] = biter.Value().([4]string)
		i++
	}

	aiter := ds.asks.Iterator()
	aiter.Begin()
	i = 0
	for aiter.Next() {
		ret.Asks[i] = aiter.Value().([4]string)
		i++
	}

	lb := ds.bids.Size()
	la := ds.asks.Size()
	fields := []string{}
	for i := 0; i < 25; i++ {
		if i < lb {
			e := ret.Bids[i]
			fields = append(fields, e[0], e[1])
		}

		if i < la {
			e := ret.Asks[i]
			fields = append(fields, e[0], e[1])
		}
	}

	raw := strings.Join(fields, ":")
	cs := crc32.ChecksumIEEE([]byte(raw))
	ret.CalcChecskum = int32(cs)
	return ret
}

func updateBook(dst *btree.Tree, elems [][4]string) {
	for _, e := range elems {
		key := e[0]
		val := e[1]

		if val == "0" {
			if _, ok := dst.Get(key); ok {
				dst.Remove(key)
			}
		} else {
			dst.Put(key, e)
		}
	}
}

func floatComparator(a, b interface{}) int {
	sa := a.(string)
	sb := b.(string)
	fa, _ := strconv.ParseFloat(sa, 64)
	fb, _ := strconv.ParseFloat(sb, 64)
	return utils.Float64Comparator(fa, fb)
}
