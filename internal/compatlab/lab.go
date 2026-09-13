package compatlab

// Compatibility laboratory per spec 55 (lines 1945-1989) and 56-60 (lines 1992-2110):
// Core matrix: Claude Code version × CCX version × protocol × provider × model × feature.
// Test categories: basic text, multi-turn, system, tools, parallel tools, tool results, reasoning,
// reasoning + tools, streaming, silent thinking, vision, documents, cache, context, cancellation,
// 429, 5xx, timeouts, disconnects, malformed SSE, model discovery, fallback.
// Golden fixtures (spec 56): every regression becomes fixture with Claude Code version, provider, model,
// feature, request/response fixtures, expected IR, expected output. Permanent regression corpus.

type Matrix struct {
	ClaudeCodeVersion string
	CCXVersion        string
	Protocol          string
	Provider          string
	Model             string
	Feature           string
}

type Fixture struct {
	ID        string
	Request   string
	Response  string
	ExpectedIR string
	ExpectedOutput string
}
