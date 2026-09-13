package protocol

// Anthropic adapter concrete implementation per spec 22 (lines 855-880) and spec 12 (Anthropic ingress, lines 532-566).
// Must handle: native Anthropic endpoint, beta headers, API version, request IDs, streaming, tool use, thinking, cache directives.

type AnthropicAdapterImpl struct{}

func (a AnthropicAdapterImpl) Validate() error {
	return nil
}

func (a AnthropicAdapterImpl) SerializeRequest(req interface{}) ([]byte, error) {
	return nil, nil
}

func (a AnthropicAdapterImpl) ParseResponse(data []byte) (interface{}, error) {
	return nil, nil
}

func (a AnthropicAdapterImpl) ParseStream(event string) (interface{}, error) {
	return nil, nil
}

func (a AnthropicAdapterImpl) MapError(err error) error {
	return err
}

func (a AnthropicAdapterImpl) DiscoverModels() ([]string, error) {
	return []string{"claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"}, nil
}
