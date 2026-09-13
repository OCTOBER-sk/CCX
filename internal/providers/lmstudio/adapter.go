package lmstudio

// LM Studio adapter per spec section 21 (lines 833-835):
// Use its supported Anthropic/native interface when available rather than translating through OpenAI unnecessarily.

type Adapter struct {
	Endpoint string
}

func New() *Adapter {
	return &Adapter{Endpoint: "http://localhost:1234"}
}
