package okex

import "net/url"

func FillsParam(instrumentID string, orderID string, before string, after string, limit string) url.Values {
	ret := url.Values{}
	pairs := [][]string{
		{"instrument_id", instrumentID},
		{"order_id", orderID},
		{"before", before},
		{"after", after},
		{"limit", limit},
	}

	for _, pair := range pairs {
		key := pair[0]
		val := pair[1]
		if val != "" {
			ret.Add(key, val)
		}
	}
	return ret
}
