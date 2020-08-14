package ftx

import (
	"context"
	"net/http"
)

type (
	Position struct {
		Cost       float64 `json:"cost"`
		EntryPrice float64 `json:"entryPrice"`
	}
)

func (client *RestClient) Positions(ctx context.Context) ([]Position, error) {
	var w Wrap
	var ret []Position

	w.Result = ret
	if err := client.request(ctx, http.MethodGet, "/positions", nil, nil, true, &w); err != nil {
		return nil, err
	}

	return ret, nil
}
