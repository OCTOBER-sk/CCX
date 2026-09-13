package context

// Context safeguards per spec 18 (lines 733-757) and 138 (lines 4452-4463):
// Separate: provider_context_limit, claude_code_context_limit, ccx_safe_context_limit.
// Safe limit = min(provider, client, policy) unless documented override exists.
// Policies: strict / safe / experimental.

type Policy string

const (
	Strict       Policy = "strict"
	Safe         Policy = "safe"
	Experimental Policy = "experimental"
)

type Limits struct {
	ProviderLimit   int
	ClientLimit     int
	SafeLimit       int
	Policy          Policy
}
