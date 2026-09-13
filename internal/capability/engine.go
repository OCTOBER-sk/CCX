package capability

// Capability engine per spec section 19 (lines 760-793):
// Each capability has DECLARED, PROBED, OBSERVED -> EFFECTIVE (SUPPORTED/UNSUPPORTED/GATED/UNVERIFIED/BROKEN).
// No arbitrary percentage score. Every capability result must include provenance (scope: provider/model/credential/protocol).

type State int

const (
	SUPPORTED State = iota
	UNSUPPORTED
	GATED
	UNVERIFIED
	BROKEN
)

type Capability struct {
	Name       string
	Declared   bool
	Probed     bool
	Observed   bool
	Effective  State
	Scope      CapabilityScope
	VerifiedAt string
	EvidenceID string
}

type CapabilityScope struct {
	Provider     string
	Model        string
	CredentialScope string
	Protocol     string
}
