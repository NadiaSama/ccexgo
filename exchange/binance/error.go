package binance

import "fmt"

type (
	//APIError error message
	APIError struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
)

func (ae *APIError) Error() string {
	return fmt.Sprintf("api error code:%d message:'%s'", ae.Code, ae.Message)
}

func (ae *APIError) Is(target interface{}) bool {
	_, ok := target.(*APIError)
	return ok
}
