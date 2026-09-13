package toolcompat

// Tool protocol normalization per spec 15 (lines 647-664) and fine-grained streaming (spec 302-305):
// Support JSON-schema/function tools, tool choice, parallel calls, results, errors, partial streamed args.
// Malformed args must produce controlled errors.
// vLLM parser differences must be explicit, not assumed from OpenAI compatibility.

type ToolType int

const (
	JSONSchema ToolType = iota
	Function
)

type Compatibility struct {
	ParallelSupported bool
	StreamingSupported bool
	ParserProfile string // e.g., "vllm-qwen" per spec 242-244
}
