package swap

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	IncomeType string

	Income struct {
		Symbol     string          `json:"symbol"`
		IncomeType IncomeType      `json:"incomeType"`
		Income     decimal.Decimal `json:"income"`
		Asset      string          `json:"asset"`
		Info       string          `json:"info"`
		Time       int64           `json:"time"`
		TranID     int64           `json:"tranId"`
		TradeID    string          `json:"tradeId"`
	}
)

const (
	IncomeTypeNone           IncomeType = ""
	IncomeTypeTransfer       IncomeType = "TRANSFER"
	IncomeTypeWelcomeBonus   IncomeType = "WELCOME_BONUS"
	IncomeTypeRealizedPnl    IncomeType = "REALIZED_PNL"
	IncomeTypeFundingFee     IncomeType = "FUNDING_FEE"
	IncomeTypeCommission     IncomeType = "COMMISSION"
	IncomeTypeInsuranceClear IncomeType = "INSURANCE_CLEAR"
)

const (
	IncomeEndPoint = "/fapi/v1/income"
)

func (rc *RestClient) Income(ctx context.Context, symbol string, it IncomeType, st int64, et int64, limit int) ([]Income, error) {
	values := url.Values{}
	if symbol != "" {
		values.Add("symbol", symbol)
	}

	if it != IncomeTypeNone {
		values.Add("incomeType", string(it))
	}

	if st != 0 {
		values.Add("startTime", fmt.Sprintf("%d", st))
	}

	if et != 0 {
		values.Add("endTime", fmt.Sprintf("%d", et))
	}

	if limit != 0 {
		values.Add("limit", fmt.Sprintf("%d", limit))
	}

	var ret []Income
	if err := rc.Request(ctx, http.MethodGet, IncomeEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get income fail")
	}
	return ret, nil
}

func (rc *RestClient) Finance(ctx context.Context, req *exchange.FinanceReqParam) ([]*exchange.Finance, error) {
	var (
		s   string
		typ IncomeType
	)
	if req.Symbol != nil {
		s = req.Symbol.String()
	}
	if req.Type == exchange.FinanceTypeFunding {
		typ = IncomeTypeFundingFee
	}

	incomes, err := rc.Income(ctx, s, typ, tconv.Time2Milli(req.StartTime),
		tconv.Time2Milli(req.EndTime), req.Limit)
	if err != nil {
		return nil, err
	}

	ret := []*exchange.Finance{}
	for i := range incomes {
		income := incomes[i]
		finance, err := income.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse income fail")
		}
		ret = append(ret, finance)
	}
	return ret, nil
}

func (ic *Income) Parse() (*exchange.Finance, error) {
	var (
		s   exchange.Symbol
		err error
	)
	if ic.Symbol != "" {
		s, err = ParseSymbol(ic.Symbol)
		if err != nil {
			return nil, err
		}
	}

	return &exchange.Finance{
		ID:       fmt.Sprintf("%d", ic.TranID),
		Time:     tconv.Milli2Time(ic.Time),
		Amount:   ic.Income,
		Symbol:   s,
		Currency: ic.Asset,
		Type:     ic.IncomeType.Parse(),
		Raw:      ic,
	}, nil
}

func (ic IncomeType) Parse() exchange.FinanceType {
	if ic == IncomeTypeFundingFee {
		return exchange.FinanceTypeFunding
	}
	return exchange.FinanceTypeOther
}
