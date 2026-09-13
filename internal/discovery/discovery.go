package discovery

// Model discovery per spec 9 (lines 416-445) and 95 (lines 3368-3409):
// Source precedence: manual explicit > project alias > provider discovery > provider catalog > CCX registry > stale cache.
// Each model record: canonical_id, provider_id, display_name, aliases, protocol, context, capabilities, pricing, region, availability, source, observed_at, expires_at.
// Stale model MUST never be silently presented as currently available.

type Source int

const (
	MANUAL Source = iota
	PROJECT_ALIAS
	PROVIDER_DISCOVERY
	PROVIDER_CATALOG
	CCX_REGISTRY
	STALE_CACHE
)

type ModelRecord struct {
	CanonicalID string
	ProviderID  string
	DisplayName string
	Aliases     []string
	Protocol    string
	Context     int
	Capabilities []string
	Pricing     interface{}
	Region      string
	Available   bool
	Source      Source
	ObservedAt  int64
	ExpiresAt   int64
	Stale       bool
}
