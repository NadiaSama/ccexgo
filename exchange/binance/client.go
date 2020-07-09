package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

type (
	//Binance Rest client instance
	RestClient struct {
		key        string
		secret     string
		recvWindow int
	}
)

func (rc *RestClient) signature(param string, body string, timestamp int64) string {
	raw := fmt.Sprintf("%s&recvWindow=%d&timestamp=%d%s", param, rc.recvWindow, timestamp, body)
	fmt.Printf("%s\n", raw)
	h := hmac.New(sha256.New, []byte(rc.secret))
	h.Write([]byte(raw))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func timeStamp() int64 {
	now := time.Now()
	return now.UnixNano() / 1e6
}
