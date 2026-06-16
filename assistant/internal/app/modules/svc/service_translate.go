package svc

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"assistant/internal/bootstrap/psl"
	"assistant/pkg/llm"
	"assistant/pkg/smartapi"
)

var translateSystemPrompt = `You are a professional translator. Translate the user's text.
Rules:
1. Auto-detect the source language.
2. If the source contains Chinese characters, translate to English.
3. Otherwise, translate to Chinese.
4. Output ONLY the translated text, nothing else.`

func (s *Service) TranslateClipboard() (string, error) {
	text, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	s.notify("translating...")

	cfg := psl.GetConfig().LLM
	client, err := llm.NewClient(cfg.Provider, llm.Config{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
		Extra:   cfg.Extra,
		Timeout: cfg.Timeout,
	})
	if err != nil {
		return "", fmt.Errorf("create LLM client: %w", err)
	}

	targetLang := detectTargetLang(text)

	translator := smartapi.NewTranslator(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := translator.Translate(ctx, text, targetLang)
	if err != nil {
		return "", fmt.Errorf("translate: %w", err)
	}

	translated := strings.TrimSpace(result.Translation)
	if translated == "" {
		return "", fmt.Errorf("empty translation result")
	}

	summary := fmt.Sprintf("translated (%s→%s)", result.SourceLanguage, result.TargetLanguage)
	return s.pushClipboard(translated, summary)
}

// detectTargetLang returns the target language based on the input text.
// If the text contains any CJK characters, target is English; otherwise Chinese.
func detectTargetLang(text string) string {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return "english"
		}
	}
	return "chinese"
}
