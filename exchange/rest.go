package exchange

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type (
	//RestReq base requeuest struct used by other
	RestReq struct {
		fields map[string]interface{}
	}
)

func NewRestReq() *RestReq {
	return &RestReq{
		fields: make(map[string]interface{}),
	}
}

func (rr *RestReq) AddFields(key string, val interface{}) *RestReq {
	rr.fields[key] = val
	return rr
}

func (rr *RestReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(rr.fields)
}

func (rr *RestReq) Values() (url.Values, error) {
	ret := url.Values{}
	for k, v := range rr.fields {
		switch t := v.(type) {
		case string:
			ret.Add(k, t)

		case int:
			ret.Add(k, strconv.FormatInt(int64(t), 10))

		case float64:
			ret.Add(k, strconv.FormatFloat(t, 'f', 8, 64))

		default:
			return ret, errors.Errorf("unknown val=%+v' for key=%s", v, k)
		}
	}
	return ret, nil
}
