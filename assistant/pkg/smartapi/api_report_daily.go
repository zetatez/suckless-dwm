package smartapi

import (
	"context"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type DailyReporter struct {
	engine *Engine
}

func NewDailyReporter(client llm.Client) *DailyReporter {
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
