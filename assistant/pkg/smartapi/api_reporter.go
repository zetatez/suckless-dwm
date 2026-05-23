package smartapi

import (
	"context"
	"fmt"

	"assistant/pkg/llm"
)

type ReportType string

const (
	DailyReport   ReportType = "daily"
	WeeklyReport  ReportType = "weekly"
	MonthlyReport ReportType = "monthly"
	QuarterReport ReportType = "quarterly"
	YearlyReport  ReportType = "yearly"
)

type ReportInput struct {
	ReportType ReportType `json:"report_type"`

	Author      string `json:"author"`
	Role        string `json:"role"`
	Period      string `json:"period"`
	Language    string `json:"language"`
	WorkContent string `json:"work_content"`
}

type ReportResult struct {
	FileName   string  `json:"file_name"`
	ReportType string  `json:"report_type"`
	Language   string  `json:"language"`
	Markdown   string  `json:"markdown"`
	Confidence float64 `json:"confidence"`
}

type Reporter struct {
	client llm.Client
}

func NewReporter(client llm.Client) *Reporter {
	return &Reporter{client: client}
}

func (r *Reporter) Generate(ctx context.Context, input ReportInput) (*ReportResult, error) {
	switch input.ReportType {
	case DailyReport:
		return NewDailyReporter(r.client).Generate(ctx, input)
	case WeeklyReport:
		return NewWeeklyReporter(r.client).Generate(ctx, input)
	case MonthlyReport:
		return NewMonthlyReporter(r.client).Generate(ctx, input)
	case QuarterReport:
		return NewQuarterlyReporter(r.client).Generate(ctx, input)
	case YearlyReport:
		return NewYearlyReporter(r.client).Generate(ctx, input)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", input.ReportType)
	}
}
