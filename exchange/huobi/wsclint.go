package huobi

import "github.com/NadiaSama/ccexgo/internal/rpc"

type (
	WSClient struct {
		url  string
		conn rpc.Conn
	}
)

func NewWSClient(addr string) *WSClient {
	return nil
}
