package stream

// Streaming state machine execution with full transition logic per spec section 13 (lines 571-608).
// Transitions: IDLE -> MESSAGE_STARTED -> BLOCK_STARTED -> DELTA -> BLOCK_STOPPED -> MESSAGE_DELTA -> MESSAGE_STOPPED.
// Terminal states must be protocol-valid (spec 90, 3176-3200). Never fabricate ambiguous events (spec 87, 2934-2972).
// Translated streams generate keepalive events for long thinking streams (spec 603-607, lines 603-607).

func (s State) IsTerminal() bool {
	return s == MESSAGE_STOPPED || s == STOPPING
}

func (s State) IsActive() bool {
	return s == MESSAGE_STARTED || s == BLOCK_STARTED || s == DELTA || s == MESSAGE_DELTA || s == RECOVERING || s == DEGRADED
}

// Transition performs state transition with spec-aligned rules.
func (s State) Transition(eventType string) State {
	switch s {
	case IDLE:
		if eventType == "message_start" {
			return MESSAGE_STARTED
		}
	case MESSAGE_STARTED:
		if eventType == "content_block_start" {
			return BLOCK_STARTED
		} else if eventType == "message_delta" {
			return MESSAGE_DELTA
		}
	case BLOCK_STARTED:
		if eventType == "delta" {
			return DELTA
		}
	case DELTA:
		if eventType == "content_block_stop" {
			return BLOCK_STOPPED
		}
	case BLOCK_STOPPED:
		if eventType == "message_delta" {
			return MESSAGE_DELTA
		}
	case MESSAGE_DELTA:
		if eventType == "message_stop" {
			return MESSAGE_STOPPED
		}
	case RECOVERING:
		if eventType == "message_stop" {
			return MESSAGE_STOPPED
		} else if eventType == "message_delta" {
			return MESSAGE_DELTA
		}
	case DEGRADED:
		if eventType == "message_stop" {
			return MESSAGE_STOPPED
		}
	}
	return s
}
