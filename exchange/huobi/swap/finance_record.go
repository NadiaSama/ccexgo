package swap

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	FinancialRecordRequest struct {
		contractCode string
		typ          []string
		createDate   int
		pageIndex    int
		pageSize     int
	}

	FinancialRecord struct {
		ID           int     `json:"id"`
		Symbol       string  `json:"symbol"`
		Type         int     `json:"type"`
		Amount       float64 `json:"amount"`
		TS           int64   `json:"ts"`
		ContractCode string  `json:"contract_code"`
	}

	FinancialRecordResponse struct {
		TotalPage       int               `json:"total_page"`
		CurrentPage     int               `json:"current_page"`
		TotalSize       int               `json:"total_size"`
		FinancialRecord []FinancialRecord `json:"financial_record"`
	}
)

const (
	FinancialRecordEndPoint = "/swap-api/v1/swap_financial_record"

	FinancialRecordTypeFundingIncome  = 30
	FinancialRecordTypeFundingOutCome = 31
)

func NewFinancialRecordRequest(contractCode string) *FinancialRecordRequest {
	return &FinancialRecordRequest{
		typ:          make([]string, 0),
		contractCode: contractCode,
	}
}

func (frr *FinancialRecordRequest) Serialize() ([]byte, error) {
	param := map[string]interface{}{
		"contract_code": frr.contractCode,
	}

	if len(frr.typ) != 0 {
		param["type"] = strings.Join(frr.typ, ",")
	}

	if frr.createDate != 0 {
		param["create_date"] = frr.createDate
	}

	if frr.pageIndex != 0 {
		param["page_index"] = frr.pageIndex
	}

	if frr.pageSize != 0 {
		param["page_size"] = frr.pageSize
	}

	return json.Marshal(param)
}
func (frr *FinancialRecordRequest) Type(types ...int) *FinancialRecordRequest {
	for _, t := range types {
		frr.typ = append(frr.typ, strconv.Itoa(t))
	}
	return frr
}

func (frr *FinancialRecordRequest) CreateDate(date int) *FinancialRecordRequest {
	frr.createDate = date
	return frr
}

func (frr *FinancialRecordRequest) PageIndex(idx int) *FinancialRecordRequest {
	frr.pageIndex = idx
	return frr
}

func (frr *FinancialRecordRequest) PageSize(size int) *FinancialRecordRequest {
	frr.pageSize = size
	return frr
}

func (cl *RestClient) FinancialRecord(ctx context.Context, req *FinancialRecordRequest) (*FinancialRecordResponse, error) {
	var ret FinancialRecordResponse
	if err := cl.PrivatePostReq(ctx, FinancialRecordEndPoint, req, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch financial record fail")
	}

	return &ret, nil
}

//Transform financeRecord to finacial currently only funding type is support
func (fr *FinancialRecord) Transform() (*exchange.Finance, error) {
	symbol, err := ParseSymbol(fr.ContractCode)
	if err != nil {
		return nil, err
	}

	if fr.Type != FinancialRecordTypeFundingIncome && fr.Type != FinancialRecordTypeFundingOutCome {
		return nil, errors.Errorf("unsupport type %d", fr.Type)
	}

	return &exchange.Finance{
		ID:       fmt.Sprintf("%d", fr.ID),
		Symbol:   symbol,
		Currency: fr.Symbol,
		Amount:   decimal.NewFromFloat(fr.Amount),
		Type:     exchange.FinanceTypeFunding,
		Time:     huobi.ParseTS(fr.TS),
		Raw:      fr,
	}, nil
}
