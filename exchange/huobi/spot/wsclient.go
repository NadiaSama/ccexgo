package spot

import (
	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

type (
	WSClient struct {
		*huobi.WSClient
	}
)

const (
	MBPAddr = "wss://api-aws.huobi.pro/feed"
)

func NewMBPWSClient(data chan interface{}) *WSClient {
	return &WSClient{
		WSClient: huobi.NewWSClient(MBPAddr, NewCodeC(), data),
	}
}

/*
func (ws *WSClient) Subscribe(ctx context.Context, channelds ...exchange.Channel) error {
	for i, ch := range channelds {
		param := huobi.CallParam{
			ID:  strconv.Itoa(i),
			Sub: ch.String(),
		}

		var dest huobi.Response

		if err := ws.Call(ctx, param.ID, huobi.MethodSubscibe, &param, &dest); err != nil {
			return err
		}

		if dest.Status != "ok" {
			return errors.Errorf("subscirbe error %+v", dest)
		}
	}
	return nil
}

*/
