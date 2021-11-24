package swap

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

func (rc *RestClient) Finance(ctx context.Context, params *exchange.FinanceReqParam) ([]exchange.Finance, error) {
	if params.Type != exchange.FinanceTypeFunding {
		return nil, errors.Errorf("unsupport type '%d'", params.Type)
	}

	if params.Symbol == nil {
		return nil, errors.Errorf("symbol is required")
	}

	req := NewFinancialRecordRequest(params.Symbol.String())
	req.Type(FinancialRecordTypeFundingIncome, FinancialRecordTypeFundingOutCome)
	if params.Limit != 0 {
		req = req.PageSize(params.Limit)
	}

	records, err := rc.FinancialRecord(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "fetch financial_record fail")
	}

	var ret []exchange.Finance
	for _, r := range records.FinancialRecord {
		rec, err := r.Transform()
		if err != nil {
			return nil, errors.WithMessage(err, "parse financial_record fail")
		}
		ret = append(ret, *rec)
	}
	return ret, nil
}
