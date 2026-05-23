package llm

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleAI     Role = "assistant"
)

type Message struct {
	Role        Role       `json:"role"`
	Content     string     `json:"content"`
	Name        string     `json:"name,omitempty"`
	ToolCalls   []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID  string     `json:"tool_call_id,omitempty"`
	ImageBase64 string     `json:"image_base64,omitempty"`
	ImageURL    string     `json:"image_url,omitempty"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float32
	MaxTokens   int
	TopP        float32 `json:"top_p,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	Tools       []Tool  `json:"tools,omitempty"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function FunctionTool `json:"function"`
}

type FunctionTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Schema `json:"parameters"`
}

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Enum        []any  `json:"enum,omitempty"`
}

type ChatResponse struct {
	Content   string
	Role      Role
	ToolCalls []ToolCall
	Usage     TokenUsage
	Raw       any
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (u TokenUsage) Total() int {
	return u.PromptTokens + u.CompletionTokens
}

type StreamCallback func(chunk ChatResponse)

type LegacyStreamCallback func(delta string)

type Capability uint64

const (
	CapabilityChat Capability = 1 << iota
	CapabilityStream
	CapabilityFunctionCall
	CapabilityVision
	CapabilityTools
)

type Capabilities struct {
	Supported Capability
}

func (c Capabilities) Has(cap Capability) bool {
	return c.Supported&cap != 0
}
