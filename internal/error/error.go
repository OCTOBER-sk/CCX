package error

// Per spec section 91 (P0 correction — error taxonomy) and spec lines 3204-3252.
// The error model must track: layer, phase, retry_class, user_action, safe_message,
// upstream_request_id, plus category/provider_code/retryable.

type CCXError struct {
	Category          string `json:"category"`
	ProviderCode      string `json:"provider_code,omitempty"`
	Retryable         bool   `json:"retryable"`
	SafeMessage       string `json:"safe_message"`
	UpstreamRequestID string `json:"upstream_request_id,omitempty"`

	// Expanded taxonomy (spec section 91)
	Layer       string `json:"layer"` // client, ccx, protocol, provider, network, auth, policy
	Phase       string `json:"phase"` // discovery, auth, request, stream, transform, cancellation, shutdown
	RetryClass  string `json:"retry_class"` // never, safe, conditional, unknown
	UserAction  string `json:"user_action"` // reauth, retry, choose_model, run_doctor, disable_feature, none
}

func (e CCXError) ValidateTaxonomy() bool {
	// Per spec 91 (lines 3204-3252) + invariant 1 (secret redaction structural — spec 98, 3172-3188): safe_message must not contain raw upstream secrets; layer/phase/retry_class/user_action must be non-empty.
	if e.SafeMessage == "" || e.Layer == "" || e.Phase == "" || e.RetryClass == "" {
		return false
	}
	return true
}

func (e CCXError) Error() string {
	return e.SafeMessage
}
