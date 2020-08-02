package huobi

type (
	Error struct {
		msg string
	}
)

func NewError(msg string) error {
	return &Error{msg}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Is(target interface{}) bool {
	_, ok := target.(*Error)
	return ok
}
