package smartapi

import (
	"context"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type LongTextSummarizer struct {
	engine *Engine
}

func NewLongTextSummarizer(client llm.Client) *LongTextSummarizer {
	return &LongTextSummarizer{engine: NewEngine(client)}
}

type SummarizeInput struct {
	Text       string   `json:"text"`
	Style      string   `json:"style,omitempty"`
	MaxLength  int      `json:"max_length,omitempty"`
	FocusAreas []string `json:"focus_areas,omitempty"`
}

type SummarizeResult struct {
	Summary    string   `json:"summary"`
	KeyPoints  []string `json:"key_points"`
	Language   string   `json:"language"`
	Confidence float64  `json:"confidence"`
}

func (s *LongTextSummarizer) Summarize(ctx context.Context, input SummarizeInput) (*SummarizeResult, error) {
	style := input.Style
	if style == "" {
		style = "standard"
	}

	maxLen := input.MaxLength
	if maxLen == 0 {
		maxLen = 500
	}

	prompt := prompts.BuildSummarizePrompt(input.Text, input.FocusAreas)
	systemPrompt := prompts.LongTextSummarizeSystemPrompt + prompts.BuildSummarizeTask(style, maxLen)

	return CompleteJSON[SummarizeResult](
		ctx,
		s.engine,
		prompt,
		systemPrompt,
		0.3,
		2048,
	)
}
