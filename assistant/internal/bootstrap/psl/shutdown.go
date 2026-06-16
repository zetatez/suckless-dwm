package psl

import (
	"context"
	"sync"
)

var (
	cleanupFuncs []func(context.Context)
	cleanupMu    sync.Mutex
)

func RegisterCleanup(fn func(context.Context)) {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	cleanupFuncs = append(cleanupFuncs, fn)
}

func ShutdownAll(ctx context.Context) {
	cleanupMu.Lock()
	funcs := cleanupFuncs
	cleanupFuncs = nil
	cleanupMu.Unlock()

	for _, fn := range funcs {
		fn(ctx)
	}
}
