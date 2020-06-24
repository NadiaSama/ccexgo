package exchange

import (
	"fmt"
)

type (
	//ErrBadArg means func argument is incorrect
	ErrBadArg struct {
		Arg interface{}
		Msg string
	}

	//ErrBadResp means exchange response message error
	ErrBadExResp struct {
		Err error
	}
)

func NewBadArg(msg string, arg interface{}) error {
	return &ErrBadArg{
		Arg: arg,
		Msg: msg,
	}
}

func (eba *ErrBadArg) Error() string {
	return fmt.Sprintf("bad arg %s %v", eba.Msg, eba.Arg)
}

func (eba *ErrBadArg) Is(target error) bool {
	_, ok := target.(*ErrBadArg)
	return ok
}

func NewBadExResp(err error) error {
	return &ErrBadExResp{Err: err}
}

func (ebe *ErrBadExResp) Error() string {
	return ebe.Err.Error()
}

func (ebe *ErrBadExResp) Is(target error) bool {
	_, ok := target.(*ErrBadExResp)
	return ok
}
