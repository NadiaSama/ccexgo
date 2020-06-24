package rpc

import "fmt"

type (
	//MsgError means the Msg format is not support by
	//the Codec and encode/decode stream data fail
	MsgError struct {
		//Msg is the message return from/send to server
		Msg []byte
		Err error
	}
	//StreamError means there is something wrong with Stream
	//underlying transport layer. StreamError will make conn
	//being closed
	StreamError struct {
		Err error
	}
)

func NewMsgError(msg []byte, err error) *MsgError {
	return &MsgError{
		Msg: msg,
		Err: err,
	}
}

func (me *MsgError) Is(err error) bool {
	_, ok := err.(*MsgError)
	return ok
}

func (me *MsgError) Error() string {
	return fmt.Sprintf("bad message: %s, %s", string(me.Msg),
		me.Err.Error())
}

func NewStreamError(err error) *StreamError {
	return &StreamError{
		Err: err,
	}
}

func (ce *StreamError) Is(err error) bool {
	_, ok := err.(*StreamError)
	return ok
}

func (ce *StreamError) Error() string {
	return fmt.Sprintf("stream error: %s", ce.Err.Error())
}
