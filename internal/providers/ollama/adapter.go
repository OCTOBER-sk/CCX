package ollama

// Ollama adapter per spec section 21 (lines 825-832):
// Do not blindly forward unsupported Claude-specific probes.
// Must recognize local runtime version/capability differences.
// Must avoid repeated unsupported token-count/probe requests.

type Adapter struct {
	Endpoint string
}

func New() *Adapter {
	return &Adapter{Endpoint: "http://localhost:11434"}
}
