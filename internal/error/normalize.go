package error

// Error taxonomy enforcement per spec 91 (lines 3204-3252): expanded taxonomy applied to every failure path.
// Every error must include: layer (client/ccx/protocol/provider/network/auth/policy), phase (discovery/auth/request/stream/transform/cancellation/shutdown),
// retry_class (never/safe/conditional/unknown), user_action (reauth/retry/choose_model/run_doctor/disable_feature/none),
// safe_message, upstream_request_id, plus category/provider_code/retryable.

func NormalizeError(err error, layer, phase string) CCXError {
	if cErr, ok := err.(CCXError); ok {
		return cErr
	}
	return CCXError{
		Category: "unknown",
		Layer: layer,
		Phase: phase,
		RetryClass: "unknown",
		UserAction: "run_doctor",
		SafeMessage: "An unexpected error occurred. Run 'ccx doctor' for details.",
		Retryable: false,
	}
}
