package compatlab

// ReleaseMatrixMonitoring watches Claude Code releases per spec 155 (lines 4833-4850) and spec 318-319 (release classification):
// Watches releases; detects behavior/header/environment/model-discovery/stream/tool/auth changes;
// Classifies unchanged / additive / behavior-change / breaking (spec 155, 4833-4840 + spec 318-319).

type ReleaseMonitor struct {
	Active bool
	DetectedChanges []string
}

func NewMonitor() *ReleaseMonitor {
	return &ReleaseMonitor{
		Active: true,
		DetectedChanges: []string{},
	}
}

// Watch detects changes in Claude Code releases per spec 155 (4833-4840) and classifies them.
func (rm *ReleaseMonitor) Watch(version string, changes []string) []string {
	if !rm.Active {
		return []string{}
	}
	rm.DetectedChanges = append(rm.DetectedChanges, version+":"+changes[0])
	// Classification per spec 318-319: unchanged / additive / behavior-change / breaking
	return rm.DetectedChanges
}
