package ir

// Universal Internal Representation (Universal Semantic IR) per spec section 11 (lines 473-532) and section 85 (expanded content blocks, lines 3006-3025).
// The IR MUST NOT stop at basic blocks; it supports an extensible registry for future Anthropic blocks.
// Unknown block behavior: known+supported → transform; known+unsupported → policy; unknown+safe → preserve; unknown+unsafe → reject.

type Request struct {
	Model        string      `json:"model"`
	System       string      `json:"system,omitempty"`
	Messages     []Message  `json:"messages"`
	Tools        []Tool     `json:"tools,omitempty"`
	ToolChoice   *ToolChoice `json:"tool_choice,omitempty"`
	Thinking     *ThinkingConfig `json:"thinking,omitempty"`
	MaxTokens    int         `json:"max_tokens,omitempty"`
	Temperature  float64     `json:"temperature,omitempty"`
	Stream       bool        `json:"stream,omitempty"`
	CacheControl []CacheControl `json:"cache_control,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Text              string             `json:"text,omitempty"`
	Image             *ImageBlock        `json:"image,omitempty"`
	Document          *DocumentBlock     `json:"document,omitempty"`
	ToolUse           *ToolUseBlock      `json:"tool_use,omitempty"`
	ToolResult        *ToolResultBlock   `json:"tool_result,omitempty"`
	Thinking          *ThinkingBlock     `json:"thinking,omitempty"`
	RedactedThinking  *RedactedThinking  `json:"redacted_thinking,omitempty"`
	ServerToolUse     *ServerToolUse     `json:"server_tool_use,omitempty"`
	WebSearchToolResult *WebSearchToolResult `json:"web_search_tool_result,omitempty"`
	SearchResult      *SearchResult      `json:"search_result,omitempty"`
	Citations         []Citation         `json:"citations,omitempty"`
	CodeExecution     *CodeExecution     `json:"code_execution,omitempty"`
	BashCodeExecution *BashCodeExecution `json:"bash_code_execution,omitempty"`
	Unknown           interface{}        `json:"-"` // preserved opaque; never silently discarded (spec 11, 532)
}

type StreamEvent struct {
	MessageStart      *MessageStartEvent `json:"message_start,omitempty"`
	ContentBlockStart *ContentBlockStartEvent `json:"content_block_start,omitempty"`
	ContentBlockDelta *ContentBlockDeltaEvent `json:"content_block_delta,omitempty"`
	ContentBlockStop  *ContentBlockStopEvent  `json:"content_block_stop,omitempty"`
	MessageDelta      *MessageDeltaEvent      `json:"message_delta,omitempty"`
	MessageStop       *MessageStopEvent       `json:"message_stop,omitempty"`
	Ping              *PingEvent              `json:"ping,omitempty"`
	Error             *StreamError            `json:"error,omitempty"`
}

type Usage struct {
	InputTokens     int `json:"input_tokens"`
	OutputTokens    int `json:"output_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens,omitempty"`
	CacheReadTokens int `json:"cache_read_tokens,omitempty"`
}

type CacheControl struct {
	Type string `json:"type"`
}

type ThinkingConfig struct {
	BudgetTokens int    `json:"budget_tokens,omitempty"`
	Type         string `json:"type,omitempty"`
}

// Minimal placeholder types for expanded content blocks (section 85)

type ImageBlock struct{}
type DocumentBlock struct{}
type ToolUseBlock struct{}
type ToolResultBlock struct{}
type ThinkingBlock struct{}
type RedactedThinking struct{}
type ServerToolUse struct{}
type WebSearchToolResult struct{}
type SearchResult struct{}
type Citation struct{}
type CodeExecution struct{}
type BashCodeExecution struct{}

type ToolChoice struct{}
type MessageStartEvent struct{}
type ContentBlockStartEvent struct{}
type ContentBlockDeltaEvent struct{}
type ContentBlockStopEvent struct{}
type MessageDeltaEvent struct{}
type MessageStopEvent struct{}
type PingEvent struct{}
type StreamError struct{}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	InputSchema interface{} `json:"input_schema,omitempty"`
}
