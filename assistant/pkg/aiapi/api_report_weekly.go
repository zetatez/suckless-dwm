package aiapi

import (
	"context"

	"assistant/pkg/llmproxy"
	"assistant/pkg/aiapi/prompts"
)

type WeeklyReporter struct {
	engine *Engine
}

func NewWeeklyReporter(client llmproxy.Client) *WeeklyReporter {
	return &WeeklyReporter{engine: NewEngine(client)}
}

func (r *WeeklyReporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	return CompleteJSON[ReportResult](
		ctx,
		r.engine,
		prompts.BuildReportContext(input.Author, input.Role, input.Period, input.Language, input.WorkContent),
		prompts.WeeklyReportSystemPrompt,
		0.4,
		2048,
	)
}
