package okex5

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	BooksEndPoint = "/api/v5/market/books"
)

func (rc *RestClient) Books(ctx context.Context, instId string, sz string) (*Depth, error) {
	var ret []Depth
	values := url.Values{}
	values.Add("instId", instId)
	if sz != "" {
		if _, err := strconv.Atoi(sz); err != nil {
			return nil, errors.WithMessagef(err, "invalid sz '%s'", sz)
		}
		values.Add("sz", sz)
	}

	if err := rc.Request(ctx, http.MethodGet, BooksEndPoint, values, nil, false, &ret); err != nil {
		return nil, err
	}

	return &ret[0], nil
}
