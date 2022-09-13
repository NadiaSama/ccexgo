package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
		startTime    int
		endTime      int
		fromID       int
		direct       string
	}

	FinancialRecord struct {
		ID           int     `json:"id"`
		Symbol       string  `json:"symbol"`
		Type         int     `json:"type"`
		Amount       float64 `json:"amount"`
		TS           int64   `json:"ts"`
		ContractCode string  `json:"contract_code"`
		QueryID      int64   `json:"query_id"`
	}

	FinancialRecordResponse struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		TS   int64             `json:"ts"`
		Data []FinancialRecord `json:"data"`
	}
)

const (
	FinancialRecordEndPoint = "/swap-api/v3/swap_financial_record"

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
		"contract": frr.contractCode,
	}

	if len(frr.typ) != 0 {
		param["type"] = strings.Join(frr.typ, ",")
	}

	if frr.startTime != 0 {
		param["start_time"] = frr.startTime
	}

	if frr.endTime != 0 {
		param["end_time"] = frr.endTime
	}

	if frr.fromID != 0 {
		param["from_id"] = frr.fromID
	}

	if frr.direct != "" {
		param["direct"] = frr.direct
	}

	return json.Marshal(param)
}
func (frr *FinancialRecordRequest) Type(types ...int) *FinancialRecordRequest {
	for _, t := range types {
		frr.typ = append(frr.typ, strconv.Itoa(t))
	}
	return frr
}

// StartTime specific record query start timestamp in milliseconds
func (frr *FinancialRecordRequest) StartTime(ts int) *FinancialRecordRequest {
	frr.startTime = ts
	return frr
}

// EndTime specific record query end timestamp in milliseconds
func (frr *FinancialRecordRequest) EndTime(ts int) *FinancialRecordRequest {
	frr.endTime = ts
	return frr
}

func (frr *FinancialRecordRequest) FromID(fromID int) *FinancialRecordRequest {
	frr.fromID = fromID
	return frr
}

func (frr *FinancialRecordRequest) Direct(direct string) *FinancialRecordRequest {
	frr.direct = direct
	return frr
}

func (cl *RestClient) FinancialRecord(ctx context.Context, req *FinancialRecordRequest) (*FinancialRecordResponse, error) {
	var ret FinancialRecordResponse

	body, err := req.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "serialzie request fail")
	}
	if err := cl.RequestWithRawResp(ctx, http.MethodPost, FinancialRecordEndPoint, nil, bytes.NewBuffer(body), true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch financial record fail")
	}

	if ret.Code != 200 {
		return nil, errors.Errorf("error response code=%d msg=%s", ret.Code, ret.Msg)
	}

	return &ret, nil
}

// Transform financeRecord to finacial currently only funding type is support
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
