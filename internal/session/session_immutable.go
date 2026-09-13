package session

// Session immutability framework per spec 102 (lines 3573-3590):
// Running session owns route_snapshot, credential_ref, model, protocol, capability_policy, transform_policy, retry_policy.
// Global config changes affect future sessions only — active session route snapshot immutable.

type ImmutableSession struct {
	Session        *Session
	RouteSnapshot  string // frozen at session start
	CredentialRef  string
	Model          string
	Protocol       string
	CapabilityPolicy string
	TransformPolicy string
	RetryPolicy     string
}

func (is *ImmutableSession) Freeze(s *Session) {
	is.Session = s
	is.RouteSnapshot = s.RouteHash
	is.CredentialRef = s.CredentialRef
	is.Model = s.Model
}

func (is *ImmutableSession) IsRouteUnchanged() bool {
	return is.RouteSnapshot == is.Session.RouteHash
}
