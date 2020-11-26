package swap

import (
	"context"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/okex"
)

func TestSwapOrder(t *testing.T) {
	key := ""
	secret := ""
	passPhrass := ""
	data := make(chan interface{}, 4)
	ctx := context.Background()

	ws := okex.NewTESTWSClient(key, secret, "", data)
	if err := ws.Auth(ctx); err == nil {
		t.Errorf("test login error fail")
	}

	//client := okex.NewTESTRestClient(key, secret, passPhrass)
	ws = okex.NewTESTWSClient(key, secret, passPhrass, data)
	if err := ws.Auth(ctx); err != nil {
		t.Fatalf("auth fail error=%s", err.Error())
	}

}
