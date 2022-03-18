package option

import (
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	DepthChannel struct {
		level int
		sym   exchange.OptionSymbol
	}
)

func NewDepthChannel(sym exchange.OptionSymbol, level int) *DepthChannel {
	return &DepthChannel{
		level: level,
		sym:   sym,
	}
}

func (dc *DepthChannel) String() string {
	return fmt.Sprintf("%s@depth%d", dc.sym.String(), dc.level)
}
