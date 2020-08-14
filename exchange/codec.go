package exchange

import "encoding/json"

type (
	//CodeC define base encode method
	CodeC struct {
	}
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

//Encode encode req with json.Marshal
func (cc *CodeC) Encode(req interface{}) ([]byte, error) {
	return json.Marshal(req)
}
