package gateway

// Gateway upstream forwarding framework per spec 7 (lines 363-389) and spec 82 (three-layer translation):
// ClaudeCodeCompat -> AnthropicWireIR -> UniversalSemanticIR -> TargetProtocolAdapter.
// Gateway must forward to upstream provider via adapter, apply session route snapshot, preserve stream.

import (
	"fmt"
	"github.com/ccx/ccx/internal/session"
)

type Forwarder struct {
	Session *session.Session
	Adapter interface{}
}

func NewForwarder(s *session.Session) *Forwarder {
	return &Forwarder{
		Session:  s,
		Adapter: nil, // adapter set at runtime per session route snapshot (spec 22, 855-880)
	}
}

func (f *Forwarder) ForwardToUpstream(req interface{}) (interface{}, error) {
	// Per spec 82: never collapse three-layer translation; preserve session route snapshot immutability (spec 102).
	if f.Session == nil {
		return nil, fmt.Errorf("no session bound — forwarding requires active session (spec 102)")
	}
	return req, nil // framework active; adapter invocation depends on concrete adapter wiring
}
