package okex5

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Trade struct {
		InstID  string    `json:"instId"`
		TradeID string    `json:"tradeId"`
		Px      string    `json:"px"`
		Sz      string    `json:"sz"`
		Side    OrderSide `json:"side"`
		Ts      string    `json:"ts"`
	}
)

const (
	TradesChannel = "trades"
)

func init() {
	parseCBMap[TradesChannel] = parseTrades
}

func NewTradesChannel(instID string) *Okex5Channel {
	return &Okex5Channel{
		Channel: TradesChannel,
		InstID:  instID,
	}
}

func parseTrades(data *wsResp) (*rpc.Notify, error) {
	var t []Trade
	if err := json.Unmarshal(data.Data, &t); err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: data.Arg.Channel,
		Params: t,
	}, nil
}
