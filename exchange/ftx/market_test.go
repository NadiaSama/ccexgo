package ftx

import (
	"context"
	"testing"
)

func TestBooks(t *testing.T) {
	client := NewRestClient("", "")
	d, err := client.Books(context.Background(), "BTC_USDT", "15")
	if err != nil {
		t.Fatalf(`Error:%v`, err)
	}
	t.Log(d)
}
