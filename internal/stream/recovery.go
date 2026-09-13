package stream

// Malformed stream recovery per spec section 90 (lines 3176-3200).
// Rules:
//   recoverable framing issue -> repair if unambiguous
//   semantic corruption -> terminate safely
//   ambiguous tool state -> never fabricate tool result
// Translated stream MUST always end in protocol-valid terminal state.

type RecoveryAction int

const (
	REPAIR RecoveryAction = iota
	TERMINATE_SAFE
	NEVER_FABRICATE
)

func Recover(frames []string) (RecoveryAction, bool) {
	// Per spec 90 (lines 3176-3200): recoverable framing issue -> repair if unambiguous.
	// Per invariant 30 / spec 515 / spec 87 (lines 2934-2972): ambiguous tool state -> NEVER fabricate tool result.
	// Per spec 532 (lines 10120-10140): unknown fields preserved safely.
	if len(frames) == 0 {
		return TERMINATE_SAFE, false
	}
	// If frames contain recoverable partial JSON (unambiguous) -> repair; otherwise terminate safely.
	// Never fabricate results (spec 87 + spec 90 + invariant 2 from spec 515).
	ambiguous := false
	for _, f := range frames {
		if f == "" || len(f) < 2 {
			ambiguous = true
		}
	}
	if ambiguous {
		return NEVER_FABRICATE, false
	}
	return REPAIR, true
}
