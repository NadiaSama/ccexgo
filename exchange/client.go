package exchange

import (
	"context"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	ConnCB func(addr string) (rpc.Conn, error)

	Client struct {
		NewConn ConnCB
		Conn    rpc.Conn
		Addr    string
		Key     string
		Secret  string
		Timeout time.Duration
		//subscribe struct will be updated via rpc Conn notify message
		Sub   map[string]interface{}
		SubMu sync.Mutex
	}
)

//NewClient got a new client instance
func NewClient(cb ConnCB, addr, key, secret string, timeout time.Duration) *Client {
	return &Client{
		NewConn: cb,
		Addr:    addr,
		Key:     key,
		Secret:  secret,
		Timeout: timeout,
		Sub:     make(map[string]interface{}),
	}
}

//Run create wsconn and start conn running loop
func (c *Client) Run(ctx context.Context) error {
	conn, err := c.NewConn(c.Addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	go c.Conn.Run(ctx, c)
	return nil
}

//Done get notify if running loop closed
func (c *Client) Done() <-chan struct{} {
	return c.Conn.Done()
}

//Error return error if running loop closed due to error
func (c *Client) Error() error {
	return c.Conn.Error()
}

//Close the running loop
func (c *Client) Close() error {
	return c.Conn.Close()
}
