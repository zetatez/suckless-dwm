package smartapi

import (
	"context"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type WeeklyReporter struct {
	engine *Engine
}

func NewWeeklyReporter(client llm.Client) *WeeklyReporter {
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
