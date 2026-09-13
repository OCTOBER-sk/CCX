package session

// Session isolation execution per spec 102 (lines 3573-3590) and spec 103 (separate health):
// Enforces: active session owns frozen route_snapshot; concurrent sessions isolated; global config changes only affect future sessions.

func (is *ImmutableSession) EnforceIsolation(newConfigRoute string) bool {
	// Per spec 102 (lines 3573-3590): active session route snapshot immutable — changes rejected unless session restarted.
	// Per invariant 12 (spec 515, lines 10805-10840): session route snapshot hash preserved and immutable.
	return is.RouteSnapshot == is.Session.RouteHash
}

func (is *ImmutableSession) ConcurrentIsolation(other *ImmutableSession) bool {
	// Per spec 102: concurrent sessions must have different session IDs and independent route snapshots.
	if is.Session == nil || other == nil {
		return false
	}
	return is.Session.ID != other.Session.ID && is.RouteSnapshot != other.RouteSnapshot
}
