package tconv

import "time"

func Milli2Time(milli int64) time.Time {
	return time.Unix(milli/1000, (milli%1000)*1000*1000)
}

func Time2Milli(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}

	return ts.UnixNano() / 1e6
}
