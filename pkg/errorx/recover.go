package errorx

import (
	"context"
	"goboot/pkg/log"
	"runtime/debug"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(logger *log.Logger, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logger.Errorf("%+v\n%s", p, debug.Stack())
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, logger *log.Logger, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logger.Errorf("%+v\n%s", p, debug.Stack())
	}
}
