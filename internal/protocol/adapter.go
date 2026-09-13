package protocol

// Protocol adapter interface per spec section 22 (lines 855-880).
// All adapters (AnthropicAdapter, OpenAIChatAdapter, OpenAIResponsesAdapter,
// GeminiAdapter, OllamaAdapter, CustomAnthropicAdapter, CustomOpenAIAdapter)
// must implement: validate, serialize_request, parse_response, parse_stream,
// map_error, discover_models.

type Adapter interface {
	Validate() error
	SerializeRequest(req interface{}) ([]byte, error)
	ParseResponse(data []byte) (interface{}, error)
	ParseStream(event string) (interface{}, error)
	MapError(err error) error
	DiscoverModels() ([]string, error)
}
