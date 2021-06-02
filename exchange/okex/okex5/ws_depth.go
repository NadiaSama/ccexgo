package okex5

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Depth struct {
		Asks     [][4]string `json:"asks"`
		Bids     [][4]string `json:"bids"`
		Ts       string      `json:"ts"`
		Checksum int         `json:"checksum"`
	}
)

const (
	DepthSnapshot = "snapshot"
	DepthUpdate   = "update"

	Books5Channel = "books5"
)

func init() {
	chs := []string{Books5Channel}

	for _, c := range chs {
		parseCBMap[c] = parseDepth
	}
}

func NewBooks5Channel(instId string) exchange.Channel {
	return &Okex5Channel{
		InstID:  instId,
		Channel: "books5",
	}
}

func parseDepth(data *wsResp) (*rpc.Notify, error) {
	var d []Depth

	if err := json.Unmarshal(data.Data, &d); err != nil {
		return nil, err
	}

	return &rpc.Notify{
		Method: data.Arg.Channel,
		Params: d,
	}, nil
}
