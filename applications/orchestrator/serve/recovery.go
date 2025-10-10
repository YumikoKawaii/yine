package serve

import (
	"fmt"
	"runtime/debug"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
)

func recoveryHandler(_ *zap.Logger) grpc_recovery.RecoveryHandlerFunc {
	return func(err interface{}) error {
		st := debug.Stack()
		limit := len(st)
		if limit > 4096 {
			limit = 4096
		}
		return fmt.Errorf("unexpected error, stacktrace (first 2048 bytes): %s, err: %+v", st[:limit], err)
	}
}
