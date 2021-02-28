package websocket

import (
	"context"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	Gen interface {
		NewConn(ctx context.Context) (Conn, error)
		Channels(ctx context.Context, oldChannel []exchange.Channel) (newChannels []exchange.Channel, nextTime time.Time, err error)
	}

	//Keeper is a struct which used to make websocket connection auto reconnect and auto update subscribe channels
	Keeper struct {
		channels []exchange.Channel
		conn     Conn
		gen      Gen
		done     chan struct{}
		ech      chan error
	}
)

func NewKeeper(gen Gen) *Keeper {
	return &Keeper{
		gen:  gen,
		done: make(chan struct{}, 0),
		ech:  make(chan error, 1),
	}
}

func (k *Keeper) Loop(ctx context.Context) {
	defer close(k.done)
	for {
		conn, err := k.gen.NewConn(ctx)
		if err == context.Canceled || err == context.DeadlineExceeded {
			return
		}
		if err != nil {
			k.pushErrorClose(err)
			continue
		}
		k.conn = conn
		k.channels = nil
		k.connLoop(ctx)
	}
}

//ECh push error when error happen
func (k *Keeper) ECh() chan error {
	return k.ech
}

func (k *Keeper) Done() chan struct{} {
	return k.done
}

func (k *Keeper) connLoop(ctx context.Context) {
	timer, err := k.updateSubscribe(ctx)
	if err != nil {
		k.pushErrorClose(err)
		return
	}
	for {
		select {
		case <-k.conn.Done():
			if err := k.conn.Error(); err != nil {
				k.pushErrorClose(err)
				return
			}
		case <-timer.C:
			timer, err = k.updateSubscribe(ctx)
			if err != nil {
				k.pushErrorClose(err)
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (k *Keeper) updateSubscribe(ctx context.Context) (*time.Timer, error) {
	channels, next, err := k.gen.Channels(ctx, k.channels)
	if err != nil {
		return nil, err
	}

	if k.channels != nil {
		if err := k.conn.UnSubscribe(ctx, k.channels...); err != nil {
			return nil, errors.WithMessage(err, "unsubscribe channel fail")
		}
	}
	if err := k.conn.Subscribe(ctx, channels...); err != nil {
		return nil, errors.WithMessage(err, "subscribe channel fail")
	}
	return time.NewTimer(time.Until(next)), nil
}

func (k *Keeper) pushErrorClose(err error) {
	if k.conn != nil {
		k.conn.Close()
	}
	k.ech <- err
}
