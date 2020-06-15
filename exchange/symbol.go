package exchange

type (
	//Symbol is used to unit different exchange markets symbol
	Symbol interface {
		String() string
	}
)