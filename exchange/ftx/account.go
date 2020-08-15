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
	var ret []Position
	if err := client.request(ctx, http.MethodGet, "/positions", nil, nil, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
