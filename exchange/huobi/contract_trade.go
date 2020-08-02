package huobi

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	Trade struct {
		Amount    float64 `json:"amount"`
		TS        int64   `json:"ts"`
		ID        int64   `json:"id"`
		Price     float64 `json:"price"`
		Direction string  `json:"direction"`
	}

	Tick struct {
		ID   int64   `json:"id"`
		TS   int64   `json:"ts"`
		Data []Trade `json:"data"`
	}
)

func parseTrades(raw json.RawMessage) ([]Trade, error) {
	var tick Tick
	if err := json.Unmarshal(raw, &tick); err != nil {
		return nil, errors.WithMessagef(err, "bad trades data %s", string(raw))
	}
	return tick.Data, nil
}
