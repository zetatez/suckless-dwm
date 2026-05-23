package psl

import (
	"context"
	"sync"
)

var (
	cleanupCancelFuncs []context.CancelFunc
	cleanupMu          sync.Mutex
)

func registerCleanup(cancel context.CancelFunc) {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	cleanupCancelFuncs = append(cleanupCancelFuncs, cancel)
}

func ShutdownAll() {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()

	for _, cancel := range cleanupCancelFuncs {
		cancel()
	}
	cleanupCancelFuncs = nil
}
