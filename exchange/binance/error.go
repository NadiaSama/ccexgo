package binance

import "fmt"

type (
	//APIError error message for /api/v3 /sapi/v3
	APIError struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}

	//WAPIError error message for wapi
	WAPIError struct {
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

func (we *WAPIError) Error() string {
	return fmt.Sprintf("wapi error message:'%s'", we.Message)
}

func (we *WAPIError) Is(target interface{}) bool {
	_, ok := target.(*WAPIError)
	return ok
}
