package fallback

// Fallback policy per spec 28 (lines 1079-1120) and 181-183 (separate budgets).
// Dimensions: credential, model, provider, profile (line 1083-1087).
// Modes: disabled, ask (default), automatic (line 1090-1096).
// Default: ask (line 1098-1102).
// Automatic fallback MUST check capability compatibility before switching (line 1104-1119).
// Fallback must preserve user intent (line 182): if active model supports vision, fallback must not silently omit it.

type Mode int

const (
	Disabled Mode = iota
	Ask
	Automatic
)

type Policy struct {
	Mode           Mode
	CheckCapabilities bool
	PreserveIntent     bool
}

func Default() Policy {
	return Policy{
		Mode: Ask,
		CheckCapabilities: true,
		PreserveIntent: true,
	}
}

// Overload protection per spec 267 (lines 6939-6947): fail fast when all candidates unhealthy.
func OverloadProtection() bool {
	return true // fail fast rather than endlessly cycling providers
}
