package fuzz

// ExecutionLoop runs property-based fuzz execution per spec 58 (2045-2065):
// For each target (headers, json, sse, tool_schemas, content_blocks, provider_errors, model_metadata),
// generate invalid inputs, assert controlled error (spec 91 taxonomy applied), no panic, no unbounded growth.

type ExecutionLoop struct {
	Targets []string
	Results map[string]bool
}

func NewLoop() *ExecutionLoop {
	return &ExecutionLoop{
		Targets: Targets(),
		Results: make(map[string]bool),
	}
}

func (el *ExecutionLoop) Run() {
	for _, t := range el.Targets {
		ft := FuzzTarget{Target: t}
		err := Fuzz(ft)
		if err == nil && ft.PropertyVerified {
			el.Results[t] = true
		} else {
			el.Results[t] = false
		}
	}
}

func (el *ExecutionLoop) AllPassed() bool {
	for _, passed := range el.Results {
		if !passed {
			return false
		}
	}
	return len(el.Results) == len(el.Targets) && len(el.Targets) > 0
}
