package binance

import (
	"time"

	"github.com/NadiaSama/ccexgo/misc/tconv"
)

func Milli2Time(ts int64) time.Time {
	return tconv.Milli2Time(ts)
}

func Time2Milli(ts time.Time) int64 {
	return tconv.Time2Milli(ts)
}
