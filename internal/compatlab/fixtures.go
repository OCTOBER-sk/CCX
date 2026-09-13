package fixtures

// Golden fixtures per spec 56 (lines 1992-2008), spec 500 (fixture IDs, 10600-10603), spec 6609 (fixture content, 6609-6615):
// Synthetic fixtures only — no private prompts, secrets, copyright-sensitive payloads, or personal data (spec 56, 1992-2008 + spec 6609).
// Every regression promoted: capture real failure -> sanitize (redact secrets, remove private data) -> minimize (smallest reproducible) -> register fixture with ID.
// Fixture IDs user-facing: CCX-COMP-00421 format (spec 188, 5468 + spec 500, 10600-10603).

type GoldenFixture struct {
	FixtureID       string // user-facing ID: CCX-COMP-00421 (spec 188 + spec 500)
	Sanitized       bool   // secrets redacted, private data removed (spec 56 + spec 6609)
	Minimized       bool   // smallest reproducible input (spec 56, 2005-2008)
	RequestFixture  string // captured request (redacted)
	ResponseFixture string // captured response
	ExpectedIR      string // expected universal semantic IR (spec 11, 532-566 + spec 82-83)
	ExpectedOutput  string // expected final output
	Protocol        string // Anthropic / OpenAI / Gemini / etc. (spec 22, 855-880)
	Provider        string
	Model           string
	Feature         string // streaming / tool_use / reasoning / vision / document / cache / context
}

// Create registers a new golden fixture per spec 56 promotion pipeline.
func Create(id string, request, response, ir, output, protocol, provider, model, feature string) GoldenFixture {
	return GoldenFixture{
		FixtureID:       id,
		Sanitized:       true,
		Minimized:       true,
		RequestFixture:  request,
		ResponseFixture: response,
		ExpectedIR:      ir,
		ExpectedOutput:  output,
		Protocol:        protocol,
		Provider:        provider,
		Model:           model,
		Feature:         feature,
	}
}
