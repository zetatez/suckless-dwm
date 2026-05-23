package smartapi

import (
	"context"
	"fmt"
	"strings"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type EmailComposer struct {
	engine *Engine
}

func NewEmailComposer(client llm.Client) *EmailComposer {
	return &EmailComposer{engine: NewEngine(client)}
}

type EmailType string

const (
	EmailTypeInquiry      EmailType = "inquiry"
	EmailTypeResponse     EmailType = "response"
	EmailTypeNotification EmailType = "notification"
	EmailTypeApology      EmailType = "apology"
	EmailTypeRequest      EmailType = "request"
	EmailTypeThankYou     EmailType = "thank_you"
	EmailTypeReminder     EmailType = "reminder"
	EmailTypeCustom       EmailType = "custom"
)

type EmailTone string

const (
	ToneFormal     EmailTone = "formal"
	ToneSemiFormal EmailTone = "semi_formal"
	ToneCasual     EmailTone = "casual"
)

type EmailInput struct {
	EmailType   EmailType `json:"email_type"`
	Tone        EmailTone `json:"tone,omitempty"`
	Language    string    `json:"language,omitempty"`
	Recipient   string    `json:"recipient"`
	Sender      string    `json:"sender"`
	Subject     string    `json:"subject,omitempty"`
	Content     string    `json:"content"`
	CC          []string  `json:"cc,omitempty"`
	BCC         []string  `json:"bcc,omitempty"`
	Attachments []string  `json:"attachments,omitempty"`
	Context     string    `json:"context,omitempty"`
}

type EmailResult struct {
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	CC         []string `json:"cc,omitempty"`
	BCC        []string `json:"bcc,omitempty"`
	Language   string   `json:"language"`
	Confidence float64  `json:"confidence"`
}

func (e *EmailComposer) Compose(ctx context.Context, input EmailInput) (*EmailResult, error) {
	tone := input.Tone
	if tone == "" {
		tone = ToneSemiFormal
	}

	lang := input.Language
	if lang == "" {
		lang = "zh-CN"
	}

	prompt := buildEmailComposePrompt(input)
	systemPrompt := fmt.Sprintf(prompts.EmailComposeSystem+`

	【本次任务】
	- 邮件类型：%s
	- 语气：%s
	`, input.EmailType, tone)

	return CompleteJSON[EmailResult](
		ctx,
		e.engine,
		prompt,
		systemPrompt,
		0.4,
		2048,
	)
}

func buildEmailComposePrompt(input EmailInput) string {
	var sb strings.Builder
	sb.WriteString("邮件基本信息：\n")
	sb.WriteString("- 收件人：" + input.Recipient + "\n")
	sb.WriteString("- 发件人：" + input.Sender + "\n")
	sb.WriteString("- 语言：" + langOrDefault(input.Language, "zh-CN") + "\n")

	if input.Subject != "" {
		sb.WriteString("- 邮件主题：" + input.Subject + "\n")
	}

	if len(input.CC) > 0 {
		sb.WriteString("- 抄送：")
		for _, cc := range input.CC {
			sb.WriteString(cc + ", ")
		}
		sb.WriteString("\n")
	}

	if len(input.Context) > 0 {
		sb.WriteString("\n背景/上下文：\n" + input.Context + "\n")
	}

	sb.WriteString("\n邮件内容/要点：\n" + input.Content + "\n")

	if len(input.Attachments) > 0 {
		sb.WriteString("\n附件：\n")
		for _, att := range input.Attachments {
			sb.WriteString("- " + att + "\n")
		}
	}

	return sb.String()
}

func langOrDefault(lang, defaultLang string) string {
	if lang == "" {
		return defaultLang
	}
	return lang
}
