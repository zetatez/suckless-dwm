package bootstrap

import (
	"context"
	"fmt"
	"time"

	"assistant/internal/app"
	"assistant/internal/bootstrap/psl"
)

func Run(ctx context.Context) error {
	if err := psl.InitConfig(); err != nil {
		return fmt.Errorf("init config failed: %w", err)
	}

	if err := psl.InitLog(); err != nil {
		return fmt.Errorf("init log failed: %w", err)
	}

	if err := psl.InitLLMClient(); err != nil {
		return fmt.Errorf("init LLM client failed: %w", err)
	}

	logger := psl.GetLogger()
	logger.Info("init log success")

	psl.RegisterCleanupLLM()
	psl.StartBackgroundTasks(ctx)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		psl.ShutdownAll(shutdownCtx)
	}()

	return app.Run(ctx)
}
