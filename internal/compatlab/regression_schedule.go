package regression

// SchedulingHook executes nightly regression scheduling per spec 155 (4833-4850) and spec 152 (test levels L0-L5):
// L0 static schema every PR; L1 wire fixtures; L2 mock provider; L3 local runtime; L4 live provider opt-in (CCX_LIVE_TESTS=1);
// L5 real Claude Code scheduled (spec 152, 4700-4705 + spec 155, 4845-4850).
// Fixture IDs user-facing: CCX-COMP-00421 (spec 188 + spec 500).

type SchedulingHook struct {
	Level int
	LiveTestsEnabled bool // L4 requires CCX_LIVE_TESTS=1 per spec 155, 4840-4845
	Scheduled bool      // L5 requires scheduled execution per spec 155, 4845-4850
}

func NewHook(level int, liveEnabled, scheduled bool) *SchedulingHook {
	return &SchedulingHook{
		Level: level,
		LiveTestsEnabled: liveEnabled,
		Scheduled: scheduled,
	}
}

func (sh *SchedulingHook) Execute() bool {
	if sh.Level < 0 || sh.Level > 5 {
		return false // spec 152 defines L0-L5 only
	}
	if sh.Level == 4 && !sh.LiveTestsEnabled {
		return false // L4 requires opt-in per spec 155
	}
	if sh.Level == 5 && !sh.Scheduled {
		return false // L5 requires scheduling per spec 155
	}
	return true
}
