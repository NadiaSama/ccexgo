package huobi

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

//ParseTS parse huobi timestmp microsecond
func ParseTS(ts int64) time.Time {
	return time.Unix(ts/1e3, ts%1e3*1e6)
}

func ParseTSStr(str string) (ret time.Time, err error) {
	var (
		ts int64
	)
	ts, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		err = errors.WithMessagef(err, "invalid timestamp ts='%s'", str)
	}

	ret = ParseTS(ts)
	return
}
