package option

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	//SubscribeReq request struct which send to binance
	wsReq struct {
		Method string      `json:"method"`
		ID     int         `json:"id"`
		Params interface{} `json:"params"`
	}

	respAccount struct {
		Total            decimal.Decimal `json:"b"`
		PosValue         decimal.Decimal `json:"m"`
		UnrealizedPNL    decimal.Decimal `json:"u"`
		OrderFrozen      decimal.Decimal `json:"o"`
		PositionFrozen   decimal.Decimal `json:"p"`
		ReduceFrozen     decimal.Decimal `json:"r"`
		MaintainceFrozen decimal.Decimal `json:"M"`
		Delta            decimal.Decimal `json:"d"`
		Theta            decimal.Decimal `json:"t"`
		Gamma            decimal.Decimal `json:"g"`
		Vega             decimal.Decimal `json:"v"`
	}

	respPosition struct {
		Symbol       string          `json:"S"`
		TotalQty     decimal.Decimal `json:"c"`
		ReducibleQty decimal.Decimal `json:"r"`
		PosValue     decimal.Decimal `json:"p"`
		AvgPrice     decimal.Decimal `json:"a"`
		RP           int             `json:"rp"`
	}

	respFilled struct {
		TradeID   string          `json:"t"`
		Price     decimal.Decimal `json:"p"`
		Qty       decimal.Decimal `json:"q"`
		TradeTime int64           `json:"T"`
		Maker     int             `json:"m"`
	}
	respOrder struct {
		Time     int64           `json:"T"`
		OrderID  string          `json:"oid"`
		Symbol   string          `json:"S"`
		Price    decimal.Decimal `json:"p"`
		Qty      decimal.Decimal `json:"q"`
		Status   int             `json:"s"`
		ExeQty   decimal.Decimal `json:"e"`
		ExeValue decimal.Decimal `json:"ec"`
		Fee      decimal.Decimal `json:"f"`
		Filled   []respFilled    `json:"fi"`
	}

	wsResp struct {
		Event    string         `json:"e"`
		TS       int64          `json:"E"`
		Symbol   string         `json:"s"`
		Ask      [][2]string    `json:"a"`
		Bid      [][2]string    `json:"b"`
		Account  []respAccount  `json:"B"`
		Position []respPosition `json:"P"`
		Order    []respOrder    `json:"o"`
		Code     int            `json:"code"`
		Desc     string         `json:"desc"`
	}

	CodeC struct {
	}
)

const (
	MethodSubscribe = "SUBSCRIBE"
	MethodPing      = "PING"
)

var (
	readerPool = &sync.Pool{}
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

func (cc *CodeC) Encode(rpc rpc.Request) ([]byte, error) {
	if rpc.Method() == MethodPing {
		return []byte("pong"), nil
	}
	i, err := strconv.Atoi(rpc.ID())
	if err != nil {
		return nil, errors.WithMessagef(err, "invalid ID='%s'", rpc.ID())
	}
	req := wsReq{
		ID:     i,
		Method: MethodSubscribe,
		Params: rpc.Params(),
	}

	return json.Marshal(req)
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	buf := bytes.NewBuffer(raw)

	var (
		reader *gzip.Reader
		err    error
	)
	if r := readerPool.Get(); r == nil {
		reader, err = gzip.NewReader(buf)
		if err != nil {
			return nil, errors.WithMessage(err, "create reader fail")
		}
	} else {
		reader = r.(*gzip.Reader)
		reader.Reset(buf)
	}

	defer func() {
		readerPool.Put(reader)
	}()

	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.WithMessage(err, "reader fail")
	}

	if len(all) == 4 && string(all) == "pong" {
		return &rpc.Result{}, nil
	}

	var resp wsResp
	if err := json.Unmarshal(all, &resp); err != nil {
		fmt.Printf("%s %s\n", string(all), err.Error())
		return nil, errors.WithMessage(err, "unmarshal json fail")
	}

	return &rpc.Notify{
		Method: "",
		Params: &resp,
	}, nil
}
