package ebus

import (
	"context"
	"fmt"
)

// Logger use to emit log
type Logger interface {
	CtxError(ctx context.Context, format string, v ...interface{})
}

var defaultLogger *_logger

type _logger struct {
}

func (log *_logger) CtxError(ctx context.Context, format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
