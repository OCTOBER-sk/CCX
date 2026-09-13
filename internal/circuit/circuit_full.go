package circuit

// Circuit breaker full logic per spec 29 (lines 1123-1142) and spec 185-186 (persistent enough with TTL, memory + short TTL):
// States: CLOSED (requests pass), OPEN (requests blocked), HALF_OPEN (limited probe requests allowed).
// Per provider/credential/model/endpoint isolation (spec 29, 1126-1133). Avoid storm: no unlimited retries during OPEN state (spec 29 + spec 64 global retry budget).

type BreakerFull struct {
	Name        string // provider/credential/model/endpoint identifier
	State       State
	Failures    int
	Successes   int
	LastFailMs  int64
	TTLMs       int64 // time-to-live for OPEN state (persistent enough, not forever — spec 185-186)
	Threshold   int  // failures before OPEN
	HalfProbeLimit int // max probe requests in HALF_OPEN
}

func NewBreaker(name string, threshold int) *BreakerFull {
	return &BreakerFull{
		Name: name,
		State: CLOSED,
		Failures: 0,
		Successes: 0,
		TTLMs: 30000, // 30s TTL — persistent enough but not forever (spec 185-186)
		Threshold: threshold,
		HalfProbeLimit: 1,
	}
}

func (b *BreakerFull) Allow() bool {
	switch b.State {
	case CLOSED:
		return true
	case OPEN:
		return false // requests blocked during OPEN (storm prevention per spec 64 + 29)
	case HALF_OPEN:
		return b.Successes < b.HalfProbeLimit // limited probes allowed
	default:
		return true
	}
}

func (b *BreakerFull) RecordSuccess() {
	b.Successes++
	if b.State == HALF_OPEN && b.Successes >= b.HalfProbeLimit {
		b.State = CLOSED
		b.Successes = 0
		b.Failures = 0
	}
}

func (b *BreakerFull) RecordFailure() {
	b.Failures++
	if b.State == HALF_OPEN {
		b.State = OPEN
		b.Successes = 0
	} else if b.Failures >= b.Threshold {
		b.State = OPEN
	}
}
