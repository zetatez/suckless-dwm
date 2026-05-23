package bootstrap

import (
	"context"
	"fmt"

	"assistant/internal/app"
	"assistant/internal/bootstrap/psl"

	_ "assistant/pkg/llm/providers/deepseek"
	_ "assistant/pkg/llm/providers/gemini"
	_ "assistant/pkg/llm/providers/glm"
	_ "assistant/pkg/llm/providers/minimax"
	_ "assistant/pkg/llm/providers/openai"
)

func Run(ctx context.Context) error {
	if err := psl.InitConfig(); err != nil {
		return fmt.Errorf("init config failed: %w", err)
	}

	if err := psl.InitLog(); err != nil {
		return fmt.Errorf("init log failed: %w", err)
	}

	logger := psl.GetLogger()
	logger.Info("init log success")

	psl.StartBackgroundTasks()

	defer func() {
		psl.ShutdownAll()
	}()

	return app.Run(ctx)
}
