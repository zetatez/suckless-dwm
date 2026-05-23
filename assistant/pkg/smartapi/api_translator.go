package smartapi

import (
	"context"
	"fmt"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type Translator struct {
	engine *Engine
}

func NewTranslator(client llm.Client) *Translator {
	return &Translator{engine: NewEngine(client)}
}

func (t *Translator) Translate(
	ctx context.Context,
	text string,
	targetLang string,
) (*TranslateResult, error) {
	prompt := fmt.Sprintf(prompts.TranslatorPromptTpl, targetLang, targetLang, text)
	return CompleteJSON[TranslateResult](
		ctx,
		t.engine,
		prompt,
		prompts.TranslatorSystem,
		0.2,
		512,
	)
}

type TranslateResult struct {
	SourceLanguage string  `json:"source_language"`
	TargetLanguage string  `json:"target_language"`
	InputType      string  `json:"input_type"`
	Translation    string  `json:"translation"`
	Confidence     float64 `json:"confidence"`
}
