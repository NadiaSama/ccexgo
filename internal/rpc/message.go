package rpc

type (
	Message interface {
	}

	Request interface {
		Message
	}

	Response interface {
		Message
	}

	//ID field for messages
	ID struct {
		Num int
	}

	//Call request message send from client to server
	Call struct {
		ID     ID
		Method string
		Params interface{}
	}

	//Result call result reply from server
	Result struct {
		ID     ID
		Method string
		Params interface{}
	}

	//Notify subscribe messages from server (kline, orders...)
	Notify struct {
		Method string
		Params interface{}
	}
)
