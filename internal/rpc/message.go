package rpc

import "encoding/json"

type (
	Message interface {
	}

	Request interface {
		Message
		ID() string
		Method() string
		Params() interface{}
	}

	Response interface {
		Message
		responseMessag() bool
	}

	//Call request message send from client to server
	Call struct {
		id     string
		method string
		params interface{}
	}

	//Error info for result code == 0 means no error
	Error struct {
		Code    int
		Message string
	}

	//Result call result reply from server
	//the Result field will be parsed via json.Unmarshal
	Result struct {
		ID     string
		Error  error
		Result json.RawMessage
	}

	//Notify subscribe messages from server (kline, orders...)
	Notify struct {
		Method string
		Params interface{}
	}
)

func NewCall(id string, method string, params interface{}) *Call {
	return &Call{
		id:     id,
		method: method,
		params: params,
	}
}

func (c *Call) ID() string {
	return c.id
}
func (c *Call) Method() string {
	return c.method
}
func (c *Call) Params() interface{} {
	return c.params
}

func (r *Result) responseMessag() bool { return true }
func (r *Notify) responseMessag() bool { return true }
