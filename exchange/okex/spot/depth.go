package spot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	depth5Raw struct {
		Asks         [][3]interface{} `json:"asks"`
		Bids         [][3]interface{} `json:"bids"`
		InstrumentID string           `json:"instrument_id"`
		Timestamp    string           `json:"timestamp"`
	}

	DepthElem struct {
		Price  float64
		Amount float64
		Orders int
	}

	//Depth5
	Depth5 struct {
		Asks   []DepthElem
		Bids   []DepthElem
		Symbol exchange.SpotSymbol
		Time   time.Time
	}

	Depth5Channel struct {
		sym exchange.SpotSymbol
	}
)

const (
	depth5Table = "spot/depth5"
)

func init() {
	okex.SubscribeCB(depth5Table, parseDepth5)
}

func NewDepth5Channel(sym exchange.SpotSymbol) exchange.Channel {
	return &Depth5Channel{
		sym: sym,
	}
}

func (dc *Depth5Channel) String() string {
	return fmt.Sprintf("%s:%s", depth5Table, dc.sym.String())
}

func parseDepth5(table string, action string, raw json.RawMessage) (*rpc.Notify, error) {
	var d depth5Raw
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}

	ts, err := okex.ParseTime(d.Timestamp)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad okex timestamp '%s'", d.Timestamp)
	}

	fields := strings.Split(d.InstrumentID, "-")
	if len(fields) != 2 {
		return nil, errors.Errorf("bad instrumetID '%s'", d.InstrumentID)
	}
	var client *okex.RestClient
	sym := client.NewSpotSymbol(fields[0], fields[1])

	processArr := func(src [][3]interface{}, dst []DepthElem) error {
		for i, v := range src {
			p, ok := v[0].(string)
			if !ok {
				return errors.Errorf("bad price value %v", v[0])
			}
			price, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return errors.WithMessagef(err, "bad price '%s'", p)
			}

			a, ok := v[1].(string)
			if !ok {
				return errors.Errorf("bad amount value %v", v[1])
			}
			amount, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return errors.WithMessagef(err, "bad amount '%s'", a)
			}

			o, ok := v[2].(int)
			if !ok {
				return errors.Errorf("bad orders value %v", v[2])
			}

			dst[i] = DepthElem{
				Amount: amount,
				Price:  price,
				Orders: o,
			}
		}
		return nil
	}
	asks := make([]DepthElem, len(d.Asks))
	if err := processArr(d.Asks, asks); err != nil {
		return nil, err
	}
	bids := make([]DepthElem, len(d.Asks))
	if err := processArr(d.Bids, bids); err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: table,
		Params: &Depth5{
			Bids:   bids,
			Asks:   asks,
			Time:   ts,
			Symbol: sym,
		},
	}, nil
}
