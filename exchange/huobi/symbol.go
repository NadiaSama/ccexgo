package huobi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	Symbol struct {
		BaseCurrency string `json:"base-currency"`
		QuoteCurreny string `json:"quote-currency"`
		Symbol       string `json:"symbol"`
	}

	SymbolResp struct {
		Status string   `json:"status"`
		Data   []Symbol `json:"data"`
	}

	FutureSymbolType int

	FutureSymbol struct {
		*exchange.BaseFutureSymbol
		typ FutureSymbolType
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

	FutureSymbolCW FutureSymbolType = iota
	FutureSymbolNW
	FutureSymbolCQ
	FutureSymbolNQ

	timeFmt = "20060102"
)

func (rc *RestClient) initSymbol(ctx context.Context) error {
	var resp SymbolResp
	if err := rc.Request(ctx, http.MethodGet, "/v1/common/symbols", nil, nil, false, &resp); err != nil {
		return err
	}

	if resp.Status != statusOK {
		ret, _ := json.Marshal(&resp)
		return newError(string(ret))
	}

	for _, sym := range resp.Data {
		rc.pair2Symbol[sym.Symbol] = rc.NewSpotSymbol(sym.BaseCurrency, sym.QuoteCurreny)
	}
	return nil
}

func (rc *RestClient) initFutureSymbol(ctx context.Context) error {
	var resp FutureSymbolResp
	if err := rc.FutureRequest(ctx, http.MethodGet, "/api/v1/contract_contract_info",
		nil, nil, false, &resp); err != nil {
		return err
	}

	if resp.Status != statusOK {
		ret, _ := json.Marshal(&resp)
		return newError(string(ret))
	}

	rc.contract2Symbol = map[string]*FutureSymbol{}
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
			typ    FutureSymbolType
		)
		if fsym.ContractType == "this_week" {
			typ = FutureSymbolCW
			suffix = "_CW"
		} else if fsym.ContractType == "next_week" {
			typ = FutureSymbolNW
			suffix = "_NW"
		} else if fsym.ContractType == "quarter" {
			typ = FutureSymbolCQ
			suffix = "_CQ"
		} else if fsym.ContractType == "next_quarter" {
			typ = FutureSymbolNQ
			suffix = "_NQ"
		} else {
			return errors.Errorf("unkown contract_type '%s'", fsym.ContractType)
		}
		sym := newFutureSymbol(fsym.Symbol, st, typ)
		rc.contract2Symbol[fmt.Sprintf("%s%s", fsym.Symbol, fsym.DeliveryDate)] = sym
		rc.contract2Symbol[fmt.Sprintf("%s%s", fsym.Symbol, suffix)] = sym
	}

	fmt.Printf("%v\n", rc.contract2Symbol)

	return nil
}

func (rc *RestClient) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToLower(base), strings.ToLower(quote)),
	}
}

func (rc *RestClient) ParseSpotSymbol(pair string) (exchange.SpotSymbol, error) {
	ret, ok := rc.pair2Symbol[pair]
	if !ok {
		return nil, errors.Errorf("unsupport pair %s", pair)
	}
	return ret, nil
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s%s", ss.Base(), ss.Quote())
}

func (rc *RestClient) GetFutureSymbols(index string) []*FutureSymbol {
	var ret []*FutureSymbol
	for k := range rc.contract2Symbol {
		v := rc.contract2Symbol[k]
		ek := fmt.Sprintf("%s%s", index, v.SettleTime().Format(timeFmt))
		if ek == k {
			ret = append(ret, v)
		}
	}
	return ret
}
func newFutureSymbol(base string, st time.Time, typ FutureSymbolType) *FutureSymbol {
	return &FutureSymbol{
		typ:              typ,
		BaseFutureSymbol: exchange.NewBaseFutureSymbol(strings.ToUpper(base), st),
	}
}

func (fs *FutureSymbol) String() string {
	return fmt.Sprintf("%s%s", fs.Index(), fs.SettleTime().Format(timeFmt))
}

//WSSub return symbol which used to subscribe trade
func (fs *FutureSymbol) WSSub() string {
	return fmt.Sprintf("%s_%s", fs.Index(), fs.typ.String())
}

func (ft FutureSymbolType) String() string {
	m := map[FutureSymbolType]string{
		FutureSymbolCW: "CW",
		FutureSymbolNW: "NW",
		FutureSymbolCQ: "CQ",
		FutureSymbolNQ: "NQ",
	}

	return m[ft]
}
