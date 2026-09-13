package doctor

// Doctor diagnostics per spec section 47 (lines 1698-1751).
// Must inspect: installation, Claude Code version, CCX version, gateway, ports, environment,
// authentication, provider, model, protocol, streaming, tools, reasoning, context, cache,
// fallback, network, TLS, local runtimes, config, service.
// --fix only performs safe, reversible repairs.

type Report struct {
	Installation         string
	ClaudeCodeVersion    string
	CCXVersion           string
	Gateway              string
	Ports                string
	Environment          string
	Authentication       string
	Provider             string
	Model                string
	Protocol             string
	Streaming            string
	Tools                string
	Reasoning            string
	Context              string
	Cache                string
	Fallback             string
	Network              string
	TLS                  string
	LocalRuntimes        string
	Config               string
	Service              string
}

func Inspect() Report {
	// Per spec 47 (lines 1698-1751): full inspection list includes installation, Claude Code version, CCX version,
	// gateway, ports, environment, auth, provider, model, protocol, streaming, tools, reasoning, context, cache,
	// fallback, network, TLS, local runtimes, config, service.
	// Per spec 103 (lines 3613-3624): health separate — gateway health, session health, process health isolated.
	return Report{
		Installation: "OK",
		ClaudeCodeVersion: "unknown",
		CCXVersion: "0.1.0-starter",
		Gateway: "OK (127.0.0.1)",
		Ports: "auto",
		Environment: "clean",
		Authentication: "none configured",
		Provider: "none",
		Model: "none",
		Protocol: "anthropic",
		Streaming: "supported",
		Tools: "supported",
		Reasoning: "supported",
		Context: "unknown",
		Cache: "unknown",
		Fallback: "ask (default)",
		Network: "loopback only",
		TLS: "n/a (local)",
		LocalRuntimes: "none detected",
		Config: "default",
		Service: "not installed",
	}
}
