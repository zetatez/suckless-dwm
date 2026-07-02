package aiapi

import (
	"context"

	"assistant/pkg/llmproxy"
	"assistant/pkg/aiapi/prompts"
)

type QuarterlyReporter struct {
	engine *Engine
}

func NewQuarterlyReporter(client llmproxy.Client) *QuarterlyReporter {
	return &QuarterlyReporter{engine: NewEngine(client)}
}

func (r *QuarterlyReporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	return CompleteJSON[ReportResult](
		ctx,
		r.engine,
		prompts.BuildReportContext(input.Author, input.Role, input.Period, input.Language, input.WorkContent),
		prompts.QuarterlyReportSystemPrompt,
		0.4,
		2048,
	)
}
