package misc

import "time"

func Milli2Time(milli int64) time.Time {
	return time.Unix(milli/1000, (milli%1000)*100)
}
