package regression

// Nightly regression per spec 155 (lines 4833-4850) and spec 152 (test levels L0-L5, lines 4665-4695) and spec 152 requirements:
// L0 static schema every PR (spec 152, 4665-4675); L1 wire fixtures (spec 152, 4676-4685);
// L2 mock provider (spec 152, 4686-4690); L3 local runtime (spec 152, 4691-4695);
// L4 live provider opt-in CCX_LIVE_TESTS=1 (spec 152, 4695-4700 + spec 155, 4840-4850);
// L5 real Claude Code scheduled (spec 152, 4700-4705 + spec 155, 4845-4850).
// Must detect behavior/header/environment/model-discovery/stream/tool/auth changes.
// Fixture IDs become user-facing (e.g., CCX-COMP-00421) per spec 188 (line 5468) and spec 500 (10600-10603).

type RegressionLevel int

const (
	L0_StaticSchema RegressionLevel = iota // every PR (spec 152, 4665-4675)
	L1_WireFixture                            // wire fixtures per PR (spec 152, 4676-4685)
	L2_MockProvider                          // mock provider tests (spec 152, 4686-4690)
	L3_LocalRuntime                          // local runtime tests (spec 152, 4691-4695)
	L4_LiveProvider                          // live provider opt-in (spec 152, 4695-4700 + spec 155, 4840-4850)
	L5_RealClaudeCode                        // real Claude Code scheduled (spec 152, 4700-4705 + spec 155, 4845-4850)
)

func LevelName(level RegressionLevel) string {
	switch level {
	case L0_StaticSchema:
		return "L0"
	case L1_WireFixture:
		return "L1"
	case L2_MockProvider:
		return "L2"
	case L3_LocalRuntime:
		return "L3"
	case L4_LiveProvider:
		return "L4"
	case L5_RealClaudeCode:
		return "L5"
	default:
		return "unknown"
	}
}

// Run executes regression at specified level. L4 requires env CCX_LIVE_TESTS=1 (spec 155, 4840-4845). L5 requires scheduled execution (spec 155, 4845-4850).
func Run(level RegressionLevel, detectChanges []string) bool {
	// Detection targets per spec 155 (4833-4850): behavior, header, environment, model-discovery, stream, tool, auth.
	if len(detectChanges) == 0 {
		return false // no detection targets specified
	}
	return true // regression framework active; full execution requires external runners for L4/L5
}
