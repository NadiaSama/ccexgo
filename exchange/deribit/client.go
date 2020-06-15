package deribit

import (
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	Client struct {
		*exchange.Client
		tokenMu     sync.Mutex
		accessToken string
		expire      time.Time
	}
)
