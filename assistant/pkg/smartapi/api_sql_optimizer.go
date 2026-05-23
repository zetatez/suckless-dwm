package smartapi

import (
	"context"
	"fmt"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type SQLOptimizer struct {
	engine *Engine
}

func NewSQLOptimizer(client llm.Client) *SQLOptimizer {
	return &SQLOptimizer{engine: NewEngine(client)}
}

func (o *SQLOptimizer) Optimize(ctx context.Context, sql string) (*SQLOptimizeResult, error) {
	prompt := fmt.Sprintf(prompts.SQLOptimizePromptTpl, sql)
	return CompleteJSON[SQLOptimizeResult](
		ctx,
		o.engine,
		prompt,
		prompts.SQLOptimizerSystem,
		0.1,
		2000,
	)
}

type SQLOptimizeResult struct {
	DatabaseType  string   `json:"database_type"`
	OriginalSQL   string   `json:"original_sql"`
	OptimizedSQL  string   `json:"optimized_sql"`
	Optimizations []string `json:"optimizations"`
	RiskLevel     string   `json:"risk_level"`
	Confidence    float64  `json:"confidence"`
}
