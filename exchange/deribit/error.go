package deribit

import (
	"fmt"
)

type (
	//JRPCError json rpc error data
	JRPCError struct {
		Code int
		Msg  string
	}
)

func NewError(code int, msg string) error {
	return &JRPCError{
		Code: code,
		Msg:  msg,
	}
}
func (je *JRPCError) Error() string {
	return fmt.Sprintf("json rpc error code: %d message: %s", je.Code, je.Msg)
}

func (js *JRPCError) Is(target error) bool {
	_, ok := target.(*JRPCError)
	return ok
}
