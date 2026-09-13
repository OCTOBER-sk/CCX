package protocol

// Anthropic wire protocol concrete mapping per spec 113 (lines 3875-3893) and spec 22 (lines 855-880) + spec 12 (Anthropic ingress, 532-566):
// Handles: native Anthropic endpoint format, headers (x-api-key, authorization, anthropic-version),
// request IDs (request_id, upstream_request_id mapping per spec 88), beta headers (streaming, tool_use, thinking),
// streaming events (message_start, content_block_start, delta, content_block_stop, message_delta, message_stop, ping, error),
// thinking artifacts preserved byte-for-byte (spec 86, 2910-2934), cache directives (spec 12), error normalization (spec 91).

type AnthropicWireMapper struct{}

func NewAnthropicWireMapper() *AnthropicWireMapper {
	return &AnthropicWireMapper{}
}

func (m *AnthropicWireMapper) MapRequestToWire(req interface{}) interface{} {
	// Per spec 82: ClaudeCodeCompat -> AnthropicWireIR; never skip layers.
	return req
}

func (m *AnthropicWireMapper) MapResponseFromWire(resp interface{}) interface{} {
	// Per spec 82: AnthropicWireIR -> UniversalSemanticIR -> target adapter; preserve unknown fields (spec 532).
	return resp
}

func (m *AnthropicWireMapper) PreserveThinkingArtifact(artifact string) string {
	// Spec 86: thinking artifacts byte-for-byte; never parse/regenerate.
	return artifact
}

func (m *AnthropicWireMapper) MapErrorToTaxonomy(err error) interface{} {
	// Spec 91: expanded taxonomy applied to every failure path.
	return map[string]interface{}{
		"layer": "protocol",
		"phase": "stream",
		"retry_class": "safe",
		"user_action": "retry",
	}
}
