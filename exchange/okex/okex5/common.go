package okex5

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

//ParseTimestamp parse okex5 api timestamp
func ParseTimestamp(ts string) (time.Time, error) {
	var t time.Time
	ret, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return t, errors.WithMessagef(err, "parse timestamp %s fail", ts)
	}

	return time.Unix(ret/1e3, ret%1e3*1e6), nil
}
