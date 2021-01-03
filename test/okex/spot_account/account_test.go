package spot_account

import (
	"context"
	"fmt"
	"testing"

	"github.com/NadiaSama/ccexgo/exchange/okex/spot"
)

func TestSpotAccount(t *testing.T) {
	ctx := context.Background()
	client := spot.NewRestClient("343126fe-8bd7-4d05-b436-4d5d4db367c1", "", "")

	accounts, err := client.FetchAccounts(ctx)
	if err != nil {
		t.Fatalf("fetch accounts fail error=%s", err.Error())
	}

	fmt.Printf("accounts=%+v\n", accounts)
}
