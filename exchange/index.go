package exchange

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type (
	//Index price
	Index struct {
		Price   float64
		Created time.Time
		Symbol  Symbol
	}

	IndexNotify Index
)

func init() {
	subRegister(reflect.TypeOf(&IndexNotify{}), indexHandler)
}

func (c *Client) Index(sym Symbol) (*Index, error) {
	c.SubMu.Lock()
	defer c.SubMu.Unlock()
	ins, ok := c.Sub[indexKey(sym)]
	if !ok {
		return nil, errors.Errorf("unkown symbol %s", sym.String())
	}
	i := ins.(*IndexNotify)
	return i.Snapshot(), nil
}

func (i *IndexNotify) Key() string {
	return indexKey(i.Symbol)
}

func (i *IndexNotify) Snapshot() *Index {
	return &Index{
		Symbol:  i.Symbol,
		Created: i.Created,
		Price:   i.Price,
	}
}

func indexHandler(ds interface{}, msg handlerMsg) interface{} {
	notify := msg.(*IndexNotify)
	return notify
}

func indexKey(sym Symbol) string {
	return fmt.Sprintf("index.%s", sym.String())
}
