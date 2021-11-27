package ctxlog

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	logKey      = struct{}{}
	ErrNoLogger = errors.New("no logger bind with ctx")
)

//GetLog get kit logger which bind with ctx. return nil if no logger set
func GetLog(ctx context.Context) log.Logger {
	ret, _ := getLogger(ctx)
	return ret
}

//SetLog bind ctx with logger
func SetLog(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

//Debug helper method. get logger from ctx and do level log via kit-level
func Debug(ctx context.Context, keyvals ...interface{}) error {
	logger, err := getLogger(ctx)
	if err != nil {
		return err
	}

	return level.Debug(logger).Log(keyvals...)
}

func Info(ctx context.Context, keyvals ...interface{}) error {
	logger, err := getLogger(ctx)
	if err != nil {
		return err
	}

	return level.Info(logger).Log(keyvals...)
}

func Warn(ctx context.Context, keyvals ...interface{}) error {
	logger, err := getLogger(ctx)
	if err != nil {
		return err
	}
	return level.Warn(logger).Log(keyvals...)
}

func Error(ctx context.Context, keyvals ...interface{}) error {
	logger, err := getLogger(ctx)
	if err != nil {
		return err
	}
	return level.Error(logger).Log(keyvals...)
}

func getLogger(ctx context.Context) (log.Logger, error) {
	logger := ctx.Value(logKey)
	if logger == nil {
		return nil, ErrNoLogger
	}

	return logger.(log.Logger), nil
}
