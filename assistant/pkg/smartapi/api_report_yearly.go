package smartapi

import (
	"context"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type YearlyReporter struct {
	engine *Engine
}

func NewYearlyReporter(client llm.Client) *YearlyReporter {
	return &YearlyReporter{engine: NewEngine(client)}
}

func (r *YearlyReporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	return CompleteJSON[ReportResult](
		ctx,
		r.engine,
		prompts.BuildReportContext(input.Author, input.Role, input.Period, input.Language, input.WorkContent),
		prompts.YearlyReportSystemPrompt,
		0.4,
		2048,
	)
}
