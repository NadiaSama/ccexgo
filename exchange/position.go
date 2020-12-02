package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	PositionSide int
	PositionMode int
	//Position info
	Position struct {
		Symbol           Symbol
		Mode             PositionMode
		Side             PositionSide
		LiquidationPrice decimal.Decimal
		AvgOpenPrice     decimal.Decimal
		CreateTime       time.Time
		Margin           decimal.Decimal
		MarginMaintRatio decimal.Decimal
		Position         decimal.Decimal
		AvailPosition    decimal.Decimal
		RealizedPNL      decimal.Decimal
		UNRealizedPNL    decimal.Decimal
		Leverage         decimal.Decimal
		Raw              interface{}
	}
)

const (
	PositionSideLong = iota
	PositionSideShort

	PositionModeFixed
	PositionModeCross
)

func (ps PositionSide) String() string {
	if PositionSideLong == ps {
		return "long"
	} else {
		return "short"
	}
}

func (pm PositionMode) String() string {
	if pm == PositionModeFixed {
		return "fixed"
	} else {
		return "crossed"
	}
}
