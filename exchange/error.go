package exchange

import "errors"

var (
	ErrUnsupport = errors.New("unsupport")
	ErrBadSymbol = errors.New("bad symbol")
)
