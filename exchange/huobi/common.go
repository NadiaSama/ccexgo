package huobi

import "time"

//ParseTS parse huobi timestmp microsecond
func ParseTS(ts int64) time.Time {
	return time.Unix(ts/1e3, ts%1e3*1e6)
}
