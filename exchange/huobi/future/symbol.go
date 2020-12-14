package future

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
)

type (
	FutureSymbol struct {
		*exchange.BaseFutureSymbol
	}

	FutureSymbolResp struct {
		Status string    `json:"status"`
		Data   []FSymbol `json:"data"`
	}

	FSymbol struct {
		Symbol         string `json:"symbol"`
		ContractCode   string `json:"contract_code"`
		DeliveryDate   string `json:"delivery_date"`
		ContractType   string `json:"contract_type"`
		ContractStatus int    `json:"contract_status"`
	}
)

const (
	ContractStatusOnline = 1
	timeFmt              = "20060102"
)

func (rc *RestClient) initFutureSymbol(ctx context.Context) error {
	var resp FutureSymbolResp
	if err := rc.Request(ctx, http.MethodGet, "/api/v1/contract_contract_info",
		nil, nil, false, &resp); err != nil {
		return err
	}

	if resp.Status != huobi.StatusOK {
		ret, _ := json.Marshal(&resp)
		return huobi.NewError(string(ret))
	}

	for _, fsym := range resp.Data {
		if fsym.ContractStatus != ContractStatusOnline {
			continue
		}

		tm := fmt.Sprintf("%s08Z", fsym.DeliveryDate)
		st, err := time.Parse("2006010215Z", tm)
		if err != nil {
			return errors.Errorf("bad delivery date %s", tm)
		}
		var (
			suffix string
			typ    exchange.FutureType
		)
		if fsym.ContractType == "this_week" {
			typ = exchange.FutureTypeCW
			suffix = "_CW"
		} else if fsym.ContractType == "next_week" {
			typ = exchange.FutureTypeNW
			suffix = "_NW"
		} else if fsym.ContractType == "quarter" {
			typ = exchange.FutureTypeCQ
			suffix = "_CQ"
		} else if fsym.ContractType == "next_quarter" {
			typ = exchange.FutureTypeNQ
			suffix = "_NQ"
		} else {
			return errors.Errorf("unkown contract_type '%s'", fsym.ContractType)
		}
		sym := newFutureSymbol(fsym.Symbol, st, typ)
		rc.futureSymbolMap[fmt.Sprintf("%s%s", fsym.Symbol, fsym.DeliveryDate)] = sym
		rc.futureSymbolMap[fmt.Sprintf("%s%s", fsym.Symbol, suffix)] = sym
	}

	return nil
}

func (rc *RestClient) GetFutureSymbols(index string) []*FutureSymbol {
	var ret []*FutureSymbol
	for k := range rc.futureSymbolMap {
		v := rc.futureSymbolMap[k]
		ek := fmt.Sprintf("%s%s", index, v.SettleTime().Format(timeFmt))
		if ek == k {
			ret = append(ret, v)
		}
	}
	return ret
}
func newFutureSymbol(base string, st time.Time, typ exchange.FutureType) *FutureSymbol {
	return &FutureSymbol{
		BaseFutureSymbol: exchange.NewBaseFutureSymbol(strings.ToUpper(base), st, typ),
	}
}

func (fs *FutureSymbol) String() string {
	return fmt.Sprintf("%s%s", fs.Index(), fs.SettleTime().Format(timeFmt))
}

//WSSub return symbol which used by websocket subscribe
func (fs *FutureSymbol) WSSub() string {
	return fmt.Sprintf("%s_%s", fs.Index(), TypeString(fs.Type()))
}

func TypeString(ft exchange.FutureType) string {
	m := map[exchange.FutureType]string{
		exchange.FutureTypeCW: "CW",
		exchange.FutureTypeNW: "NW",
		exchange.FutureTypeCQ: "CQ",
		exchange.FutureTypeNQ: "NQ",
	}

	return m[ft]
}
