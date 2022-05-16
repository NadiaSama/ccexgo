package binance

import "fmt"

type (
	APIIF interface {
		ECode() int
		EMessage() string
	}

	//APIError error message
	APIError struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
)

func (ae *APIError) ECode() int {
	return ae.Code
}

func (ae *APIError) EMessage() string {
	return ae.Message
}

func (ae *APIError) Error() string {
	return fmt.Sprintf("api error code:%d message:'%s'", ae.Code, ae.Message)
}

func (ae *APIError) Is(target interface{}) bool {
	_, ok := target.(*APIError)
	return ok
}
