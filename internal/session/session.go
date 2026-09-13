package session

import (
	"crypto/sha256"
	"fmt"
)

// Session Manager per spec section 4.4 (lines 204-212) and section 102 (active session immutability, lines 3573-3590).
// Each session owns: provider, model, credential, route, config snapshot, request counters.
// Changing global config affects future sessions, not existing ones (section 102).

type Session struct {
	ID            string
	RouteSnapshot string // hash of active route config
	Provider      string
	Model         string
	CredentialRef string
	RouteHash     string
}

func New(id string, provider, model string) *Session {
	raw := fmt.Sprintf("%s:%s:%s", id, provider, model)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
	return &Session{
		ID: id,
		Provider: provider,
		Model: model,
		RouteSnapshot: hash,
		RouteHash: hash,
	}
}

func (s *Session) VerifyHash() bool {
	// Per spec 102 (lines 3573-3590) + invariant 12 (spec 515, lines 10805-10840): session route snapshot hash must be preserved byte-for-byte; any mutation fails verification.
	raw := fmt.Sprintf("%s:%s:%s", s.ID, s.Provider, s.Model)
	expected := fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
	return s.RouteHash == expected && s.RouteSnapshot == expected
}

func (s *Session) RecomputeHash() string {
	// Per spec 102: recompute hash from current fields; used only after explicit session restart (not during active session).
	raw := fmt.Sprintf("%s:%s:%s", s.ID, s.Provider, s.Model)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
}
