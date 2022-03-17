package option

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	//Symbol implement OptionSymbol interface
	Symbol struct {
		*exchange.BaseOptionSymbol
		symbol string
	}

	//OptionInfo option info of binance option contract
	OptionInfo struct {
		ID                   int             `json:"id"`
		ContractID           int             `json:"contractId"`
		Underlying           string          `json:"underlying"`
		QuoteAsset           string          `json:"quoteAsset"`
		Symbol               string          `json:"symbol"`
		Unit                 decimal.Decimal `json:"unit"`
		MinQty               decimal.Decimal `json:"minQty"`
		MaxQty               decimal.Decimal `json:"maxQty"`
		PriceScale           int             `json:"priceScale"`
		QuantityScale        int             `json:"quantityScale"`
		Side                 string          `json:"side"`
		Leverage             decimal.Decimal `json:"leverage"`
		StrikePrice          decimal.Decimal `json:"strikePrice"`
		MakerFeeRate         decimal.Decimal `json:"makerFeeRate"`
		TakerFeeRate         decimal.Decimal `json:"takerFeeRate"`
		InitialMargin        decimal.Decimal `json:"initialMargin"`
		AutoReduceMargin     decimal.Decimal `json:"autoReduceMargin"`
		MaintenanceMargin    decimal.Decimal `json:"maintenanceMargin"`
		MinInitialMargin     decimal.Decimal `json:"minInitialMargin"`
		MinAutoReduceMargin  decimal.Decimal `json:"minAutoReduceMargin"`
		MinMaintenanceMargin decimal.Decimal `json:"minMaintenanceMargin"`
		ExpiryDate           int64           `json:"expiryDate"`
	}
)

const (
	OptionInfoEndPoint = "/vapi/v1/optionInfo"
	SideCall           = "CALL"
	SidePut            = "PUT"
)

var (
	symbolMap  = map[string]*Symbol{}
	mu         sync.Mutex
	useTestNet bool
)

func Init(ctx context.Context, testNet bool) error {
	useTestNet = testNet
	st, err := updateSymbol(ctx)
	if err != nil {
		return errors.WithMessage(err, "updateSymbol fail")
	}

	go func() {
		var (
			timer *time.Timer
		)

		if st.After(time.Now()) {
			timer = time.NewTimer(time.Until(st.Add(time.Second * 10)))
		} else {
			timer = time.NewTimer(time.Second)
		}

		for {
			select {
			case <-ctx.Done():
				return

			case <-timer.C:
				st, err = updateSymbol(ctx)
				if err != nil {
					timer = time.NewTimer(time.Second * 5)
				} else {
					timer = time.NewTimer(time.Until(st.Add(time.Second * 10)))
				}
			}
		}
	}()
	return nil
}

func ParseSymbol(symbol string) (exchange.OptionSymbol, error) {
	mu.Lock()
	defer mu.Unlock()
	if sym, ok := symbolMap[symbol]; ok {
		return sym, nil
	} else {
		return nil, errors.Errorf("unknown symbol='%s'", symbol)
	}
}

func updateSymbol(ctx context.Context) (minSettleTime time.Time, err error) {
	var rc *RestClient

	if useTestNet {
		rc = NewTestRestClient("", "")
	} else {
		rc = NewRestClient("", "")
	}

	var (
		symbols []exchange.OptionSymbol
	)
	symbols, err = rc.Symbols(ctx)
	if err != nil {
		err = errors.WithMessage(err, "fetch symbols fail")
		return
	}

	if len(symbols) == 0 {
		err = errors.Errorf("no symbols")
		return
	}

	sort.Slice(symbols, func(i, j int) bool {
		si := symbols[i]
		sj := symbols[j]

		return si.SettleTime().Before(sj.SettleTime())
	})

	nm := make(map[string]*Symbol, len(symbols))
	for _, s := range symbols {
		sym := s.(*Symbol)
		nm[s.String()] = sym
	}

	minSettleTime = symbols[0].SettleTime()

	mu.Lock()
	defer mu.Unlock()
	symbolMap = nm
	return
}

func (rc *RestClient) OptionInfo(ctx context.Context) ([]OptionInfo, error) {
	var oi []OptionInfo

	if err := rc.GetRequest(ctx, OptionInfoEndPoint, binance.NewRestReq(), false, &oi); err != nil {
		return nil, err
	}
	return oi, nil
}

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.OptionSymbol, error) {
	ois, err := rc.OptionInfo(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get option info fail")
	}

	ret := make([]exchange.OptionSymbol, len(ois))
	for i := range ois {
		oi := ois[i]
		sym, err := oi.Parse()
		if err != nil {
			return nil, errors.WithMessagef(err, "parse symbol fail symbol='%s'", oi.Symbol)
		}
		ret[i] = sym
	}
	return ret, nil
}

func (oi *OptionInfo) Parse() (*Symbol, error) {
	var typ exchange.OptionType
	if oi.Side == SideCall {
		typ = exchange.OptionTypeCall
	} else if oi.Side == SidePut {
		typ = exchange.OptionTypePut
	} else {
		return nil, errors.Errorf("unkown side '%s'", oi.Side)
	}

	bo := exchange.NewBaseOptionSymbol(oi.Underlying, time.Unix(oi.ExpiryDate/1e3, 0),
		oi.StrikePrice, typ, exchange.SymbolConfig{
			PricePrecision:  decimal.New(1, int32(oi.PriceScale)*-1),
			AmountPrecision: decimal.New(1, int32(oi.QuantityScale)*-1),
			AmountMin:       oi.MinQty,
			AmountMax:       oi.MaxQty,
		}, oi)

	return &Symbol{
		BaseOptionSymbol: bo,
		symbol:           oi.Symbol,
	}, nil
}

func (s *Symbol) String() string {
	return s.symbol
}
