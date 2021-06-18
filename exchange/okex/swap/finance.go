package swap

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
)

func (rc *RestClient) Ledgers(ctx context.Context, instrumentID string, before, after, limit, typ string) ([]okex.Ledger, error) {
	endPoint := fmt.Sprintf("/api/swap/v3/accounts/%s/ledger", instrumentID)

	ret, err := okex.FetchLedgers(ctx, rc, endPoint, before, after, limit, typ)
	return ret, err
}

func (rc *RestClient) Finance(ctx context.Context, req *exchange.FinanceReqParam) ([]exchange.Finance, error) {
	var symbol string
	if req.Symbol != nil {
		symbol = req.Symbol.String()
	}

	var typ string
	if req.Type == exchange.FinanceTypeFunding {
		typ = "14"
	}
	ledgers, err := rc.Ledgers(ctx, symbol, req.StartID, req.EndID, strconv.Itoa(req.Limit), typ)
	if err != nil {
		return nil, err
	}
	ret := []exchange.Finance{}

	for i := range ledgers {
		ledger := ledgers[i]
		f, err := ledger.Parse(parseSymbol)
		if err != nil {
			return nil, errors.WithMessage(err, "parse ledger fail")
		}
		ret = append(ret, *f)
	}
	return ret, nil
}

func parseSymbol(sym string) (exchange.Symbol, error) {
	s, err := ParseSymbol(sym)
	if err != nil {
		return nil, err
	}

	return s.(exchange.Symbol), err
}
