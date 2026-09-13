package profiles

// Profile execution logic per spec 39 (lines 1356-1383) and spec 288-289 (requirements/candidates/fallback/probe):
// Profile is a route policy: selects candidates matching requirements, applies fallback when required capabilities unsupported,
// runs probe policy to verify capability before routing.

func (p Profile) MatchCandidate(candidate string) bool {
	for _, c := range p.Candidates {
		if c == candidate {
			return true
		}
	}
	return false
}

func (p Profile) HasRequirements() bool {
	return p.Requirements.Tools != "" || p.Requirements.Reasoning != "" || p.Requirements.Vision != "" || p.Requirements.ContextMin > 0
}

func (p Profile) Validate() bool {
	if p.Name == "" || len(p.Candidates) == 0 {
		return false // spec 39: profile requires name and at least one candidate
	}
	return true
}
