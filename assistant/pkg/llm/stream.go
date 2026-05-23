package llm

type StreamChunk struct {
	Content   string
	Role      Role
	Done      bool
	Usage     TokenUsage
	ToolCalls []ToolCall
}

func (s StreamChunk) IsEmpty() bool {
	return s.Content == "" && !s.Done && len(s.ToolCalls) == 0
}
