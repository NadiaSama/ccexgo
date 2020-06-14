package rpc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type (
	Handler interface {
		Handle(ctx context.Context, n *Notify)
	}

	//Conn a connection bettwen the client and server
	Conn interface {
		//Call send request from client to server. if r is not nil the
		//Call will be waiting for the server return and set Result
		Call(ctx context.Context, method string, params interface{}, r *Result) error
		//Run start a gorotuine loop for notify message from server
		//and call handler for each message
		Run(ctx context.Context, handler Handler)
		//Done being closed if running loop was done
		Done() <-chan struct{}
		//Error return err if running loop stop due to error
		Error() error
		//Close stop the running loop
		Close() error
	}

	connection struct {
		stream    Stream
		pending   map[ID]chan *rpcCall
		done      chan struct{}
		err       atomic.Value
		seq       int64
		streamMu  sync.Mutex
		pendingMu sync.Mutex
	}

	rpcCall struct {
		err    error
		result *Result
	}
)

var (
	ErrClear = errors.New("handleMessages closed")
)

func NewConn(stream Stream) Conn {
	return &connection{
		stream:  stream,
		pending: make(map[ID]chan *rpcCall),
		done:    make(chan struct{}),
	}
}

func (c *connection) Call(ctx context.Context, method string, params interface{}, r *Result) error {
	var err error
	var rchan chan *rpcCall
	call := NewCall(atomic.AddInt64(&c.seq, 1), method, params)

	if r != nil {
		rchan = make(chan *rpcCall, 1)
		c.pendingMu.Lock()
		c.pending[call.id] = rchan
		c.pendingMu.Unlock()

		defer func() {
			c.pendingMu.Lock()
			delete(c.pending, call.id)
			c.pendingMu.Unlock()
		}()
	}

	if err = c.write(call); err != nil {
		if errors.Is(err, &StreamError{}) {
			c.fail(err)
		}
		return err
	}

	if r == nil {
		return nil
	}

	select {
	case rc, ok := <-rchan:
		if !ok {
			//handleMessage quit
			return ErrClear
		}
		if rc.err != nil {
			return rc.err
		}
		*r = *rc.result
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connection) Run(ctx context.Context, handler Handler) {
	go c.handleMessages(ctx, handler)
}

func (c *connection) Done() <-chan struct{} {
	return c.done
}

func (c *connection) Error() error {
	if err := c.err.Load(); err != nil {
		return err.(error)
	}
	return nil
}

func (c *connection) Close() error {
	return c.stream.Close()
}

func (c *connection) write(call *Call) error {
	c.streamMu.Lock()
	defer c.streamMu.Unlock()
	return c.stream.Write(call)
}

//save err value and close stream
func (c *connection) fail(err error) {
	c.err.Store(err)
	c.stream.Close()
}

//close all pending rpcCall channel
func (c *connection) clear() {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	for _, ch := range c.pending {
		close(ch)
	}
}

func (c *connection) handleMessages(ctx context.Context, handler Handler) {
	defer close(c.done)
	defer c.clear()
	for {
		response, err := c.stream.Read()
		if err != nil {
			if errors.Is(err, &StreamError{}) {
				c.fail(err)
				return
			}
		}

		switch msg := response.(type) {
		case *Result:
			c.pendingMu.Lock()
			rchan, ok := c.pending[msg.ID]
			c.pendingMu.Unlock()
			if ok {
				rchan <- &rpcCall{
					err:    nil,
					result: msg,
				}
			}
			break

		case *Notify:
			handler.Handle(ctx, msg)
			break
		}
	}
}
