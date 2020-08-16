package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
)

type (
	testStream struct {
		result []string
		lastID string
		closed bool
		cond   *sync.Cond
		wait   bool
	}
)

var (
	testErr = errors.New("test error")
)

func (ts *testStream) Read() (Response, error) {
	//wait Write
	ts.cond.L.Lock()
	defer ts.cond.L.Unlock()
	defer func() { ts.wait = true }()
	for ts.wait {
		ts.cond.Wait()
	}

	if len(ts.result) == 0 {
		return nil, testErr
	}

	val := ts.result[0]
	ts.result = ts.result[1:]
	return &Result{ID: ts.lastID, Result: json.RawMessage(val)}, nil
}

func (ts *testStream) Write(req Request) error {
	if ts.closed {
		return NewStreamError(errors.New("closed"))
	}
	call := req.(*Call)
	ts.lastID = call.id
	ts.wait = false
	ts.cond.Broadcast()
	return nil
}
func (ts *testStream) Close() error {
	ts.closed = true
	return nil
}

func TestCall(t *testing.T) {
	stream := &testStream{
		result: []string{"[1]", "[2]", "[3]", "[4]"},
		closed: false,
		cond:   sync.NewCond(&sync.Mutex{}),
		wait:   true,
	}

	conn := NewConn(stream)
	ctx := context.Background()
	conn.Run(ctx, nil)
	var arr []int
	for i := 1; i < 5; i++ {
		conn.Call(ctx, strconv.Itoa(i), "", nil, &arr)
		if arr[0] != i {
			t.Errorf("bad value %v", arr)
		}
	}
	//read error ctx will timeout
	ctx2, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	conn.Call(ctx2, "", "", nil, &arr)
	conn.Close()
	if err := conn.Call(ctx2, "", "", nil, &arr); !errors.Is(err, &StreamError{}) {
		t.Errorf("bad expect error %v", err)
	}
	c := conn.(*connection)
	if len(c.pending) != 0 {
		t.Errorf("bad state for connection %v", c.pending)
	}
}

type (
	testStreamH struct {
		testStream
	}
)

func (tsh *testStreamH) Read() (Response, error) {
	return nil, NewStreamError(errors.New("test stream handler quit"))
}

func TestHadleMessagQuit(t *testing.T) {
	stream := &testStreamH{
		testStream: testStream{
			result: []string{},
			closed: false,
			cond:   sync.NewCond(&sync.Mutex{}),
			wait:   true,
		},
	}

	conn := NewConn(stream)
	result := &Result{}
	ctx := context.Background()
	conn.Run(ctx, nil)

	if err := conn.Call(ctx, "id1", "", nil, result); err != ErrClear {
		t.Errorf("bad error %v", err)
	}
}

func TestClose(t *testing.T) {
	stream := &testStreamH{
		testStream: testStream{
			result: []string{},
			closed: false,
			cond:   sync.NewCond(&sync.Mutex{}),
			wait:   true,
		},
	}

	conn := NewConn(stream)
	ctx, cancel := context.WithCancel(context.Background())
	conn.Run(ctx, nil)

	select {
	case <-time.After(time.Second * 2):
		cancel()
	}

	select {
	case <-conn.Done():
	case <-time.After(time.Second):
		t.Fatalf("context close timeout")
	}
}
