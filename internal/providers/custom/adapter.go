package custom

// Custom Anthropic endpoint adapter per spec section 24 (lines 927-959) and 112 (declarative config):
// Most custom providers should require: provider id, protocol, base_url, auth, models.
// Auto protocol detection must be conservative; if ambiguous, ask (line 958).

type Adapter struct {
	ID       string
	Protocol string // anthropic, openai-chat, openai-responses, gemini, auto
	BaseURL  string
	AuthType string // api_key, bearer, custom_header, none
}

func New(id string, protocol string, baseURL string, authType string) *Adapter {
	return &Adapter{
		ID: id,
		Protocol: protocol,
		BaseURL: baseURL,
		AuthType: authType,
	}
}
