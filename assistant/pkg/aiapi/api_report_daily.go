package aiapi

import (
	"context"

	"assistant/pkg/llmproxy"
	"assistant/pkg/aiapi/prompts"
)

type DailyReporter struct {
	engine *Engine
}

func NewDailyReporter(client llmproxy.Client) *DailyReporter {
	return &DailyReporter{engine: NewEngine(client)}
}

func (r *DailyReporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	return CompleteJSON[ReportResult](
		ctx,
		r.engine,
		prompts.BuildReportContext(input.Author, input.Role, input.Period, input.Language, input.WorkContent),
		prompts.DailyReportSystemPrompt,
		0.4,
		2048,
	)
}
