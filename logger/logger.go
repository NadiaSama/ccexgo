package logger

import "context"

type (
	Logger interface {
		Log(keyvals ...interface{}) error
	}

	keyType struct {
	}
)

var (
	loggerKey = keyType{}
)

//Bind return a copy of ctx with speicifc logger bind
func Bind(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

//Extract the logger which bind with ctx. return nil if no logger bind
func Extract(ctx context.Context) Logger {
	ret := ctx.Value(loggerKey)
	if ret == nil {
		return nil
	}
	return ret.(Logger)
}
