package exchange

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	//CodeC define base encode method
	CodeC struct {
	}
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

//Encode encode req with json.Marshal
func (cc *CodeC) Encode(req rpc.Request) ([]byte, error) {
	return json.Marshal(req.Params())
}
