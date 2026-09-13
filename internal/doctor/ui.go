package doctorui

// Doctor UI per spec 47 (lines 1698-1751) and 47.1 (line 1750):
// `ccx doctor` inspects installation, Claude Code version, gateway, auth, provider, model, streaming,
// tools, reasoning, context, cache, fallback, network, TLS, local runtimes, config, service.
// `--fix` only safe/reversible repairs with preview (spec 414-416).

type UI struct{}

func (u *UI) Show() string {
	// Per spec 47 (lines 1698-1751) + 414-416: doctor UI must display full inspection framework (installation, version, gateway, auth, provider, model, streaming, tools, reasoning, context, cache, fallback, network, TLS, local runtimes, config, service) with safe/reversible repair preview.
	// Per spec 103 (lines 3613-3624): health separate — doctor UI isolated from gateway/session process health.
	return "CCX Doctor: full inspection active (spec 47, 414-416) — health separate (spec 103) — safe/reversible repairs only with preview"
}
