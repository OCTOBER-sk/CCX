package tui

// TUI per spec 43 (lines 1537-1588):
// Theme: background graphite/near-black, primary white, accent electric cyan.
// Brand communicates routing, connection, infrastructure, reliability.
// Main dashboard shows profile, provider, model, protocol, capabilities, sessions, fallback, gateway.
// Must work at 80x24 minimum (spec 159), keyboard-only (spec 158), with persistent route pill.

type Theme struct {
	Background string
	Primary    string
	Accent     string
	Success    string
	Warning    string
	Danger     string
	Muted      string
}

var DefaultTheme = Theme{
	Background: "#1a1a1a",
	Primary:    "#ffffff",
	Accent:     "#00d4ff",
	Success:    "#00aa00",
	Warning:    "#ffaa00",
	Danger:     "#aa0000",
	Muted:      "#888888",
}
