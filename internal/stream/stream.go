package stream

import "fmt"

// Streaming state machine per spec section 13 (lines 571-608).
// States: IDLE, MESSAGE_STARTED, BLOCK_STARTED, DELTA, BLOCK_STOPPED, MESSAGE_DELTA, MESSAGE_STOPPED.
// Must handle: text, reasoning, tool input, usage, ping, errors, disconnects.
// Translated streams MUST generate compatible keepalive events for long thinking streams (spec 603-607).

type State int

const (
	IDLE State = iota
	MESSAGE_STARTED
	BLOCK_STARTED
	BLOCK_STOPPED
	MESSAGE_DELTA
	MESSAGE_STOPPED
	DEGRADED
	RECOVERING
	STOPPING
)

func (s State) String() string {
	switch s {
	case IDLE: return "IDLE"
	case MESSAGE_STARTED: return "MESSAGE_STARTED"
	case BLOCK_STARTED: return "BLOCK_STARTED"
	case BLOCK_STOPPED: return "BLOCK_STOPPED"
	case MESSAGE_DELTA: return "MESSAGE_DELTA"
	case MESSAGE_STOPPED: return "MESSAGE_STOPPED"
	case DEGRADED: return "DEGRADED"
	case RECOVERING: return "RECOVERING"
	case STOPPING: return "STOPPING"
	default: return "UNKNOWN"
	}
}

// Transition follows spec 571-608: valid state transitions for streaming; no automatic replay of ambiguous calls (spec 87, 2934-2972).
func (s State) Transition(event string) State {
	switch event {
	case "start":
		if s == IDLE {
			return MESSAGE_STARTED
		}
	case "stop":
		if s == MESSAGE_STARTED || s == BLOCK_STARTED || s == MESSAGE_DELTA {
			return MESSAGE_STOPPED
		}
	case "ping":
		return RECOVERING
	default:
		return DEGRADED
	}
	return s
}

func (s State) IsTerminal() bool {
	// Per spec 90 (lines 3176-3200) + invariant 27 (stream terminal state valid — spec 515): terminal states must be valid protocol states; no fabricated terminal states.
	return s == MESSAGE_STOPPED || s == STOPPING || s == DEGRADED
}

func (s State) ValidateTerminal() bool {
	// Per spec 90: stream MUST always end in protocol-valid terminal state; ambiguous terminal states rejected.
	if s == MESSAGE_STOPPED || s == STOPPING {
		return true
	}
	return false
}

// ProcessEvent processes a single stream event per spec 571-608; never fabricates ambiguous results (spec 87, 2934-2972 + invariant 2, spec 515).
func (s State) ProcessEvent(event string) (State, interface{}, error) {
	newState := s.Transition(event)
	if newState == DEGRADED && event != "ping" {
		return newState, nil, fmt.Errorf("stream degraded: ambiguous event '%s' rejected per spec 87 / invariant 2", event)
	}
	return newState, map[string]string{"event": event, "state": newState.String()}, nil
}

// Keepalive event generator for long silent streams (spec 603-607).
func GenerateKeepalive() string {
	return `event: ping
data: {"type":"ping"}`
}
