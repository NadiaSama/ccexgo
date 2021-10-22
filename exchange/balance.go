package exchange

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Balance struct {
		Currency string
		Total    decimal.Decimal
		Free     decimal.Decimal
		Frozen   decimal.Decimal
	}

	Balances struct {
		Balances map[string]*Balance
		Raw      interface{}
	}
)

func NewBalances() *Balances {
	return &Balances{
		Balances: make(map[string]*Balance),
	}
}

//Get return balance for specific token(upper case)
func (b *Balances) Get(currency string) (*Balance, error) {
	r, ok := b.Balances[currency]
	if !ok {
		return nil, errors.Errorf("no balance for '%s'", currency)
	}

	return r, nil
}
