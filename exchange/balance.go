package exchange

import (
	"strings"

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

//CurrencyFormat common method which used to transfer currency to uniq format
func CurrencyFormat(input string) string {
	return strings.ToUpper(input)
}

func NewBalances() *Balances {
	return &Balances{
		Balances: make(map[string]*Balance),
	}
}

//Add specific balance
func (b *Balances) Add(balance *Balance) {
	new := *balance
	new.Currency = CurrencyFormat(balance.Currency)
	b.Balances[balance.Currency] = &new
}

//Get return balance for specific token(upper case)
func (b *Balances) Get(currency string) (*Balance, error) {
	r, ok := b.Balances[CurrencyFormat(currency)]
	if !ok {
		return nil, errors.Errorf("no balance for '%s'", currency)
	}

	return r, nil
}
