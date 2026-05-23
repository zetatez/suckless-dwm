package smartapi

import (
	"context"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type MonthlyReporter struct {
	engine *Engine
}

func NewMonthlyReporter(client llm.Client) *MonthlyReporter {
	return &MonthlyReporter{engine: NewEngine(client)}
}

func (r *MonthlyReporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	return CompleteJSON[ReportResult](
		ctx,
		r.engine,
		prompts.BuildReportContext(input.Author, input.Role, input.Period, input.Language, input.WorkContent),
		prompts.MonthlyReportSystemPrompt,
		0.4,
		2048,
	)
}
