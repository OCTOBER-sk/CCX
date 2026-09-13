package profiles

// Profiles per spec 39 (lines 1356-1383): profile is a route policy (not model alias).
// V1 profile: requirements + candidates + fallback policy + probe policy (spec 288-289).
// No complex inheritance in V1 (spec 289).

type Profile struct {
	Name         string
	Requirements Requirements
	Candidates   []string
	Fallback     string
	ProbePolicy  string
}

type Requirements struct {
	Tools      string // required, preferred, optional, forbidden
	Reasoning  string
	Vision     string
	ContextMin int
}
