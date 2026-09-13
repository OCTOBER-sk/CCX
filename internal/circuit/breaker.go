package circuit

// Circuit breaker per spec 29 (lines 1123-1142) and 185-186 (persistent enough with TTL):
// States: CLOSED, OPEN, HALF_OPEN (line 1135-1141).
// Track per provider/credential/model/endpoint (line 1126-1133).
// Avoid hammering failing endpoints. State must not persist forever (line 185-186: memory + short TTL).

type State int

const (
	CLOSED State = iota
	OPEN
	HALF_OPEN
)

type Breaker struct {
	State      State
	Failures   int
	LastFail   int64
	TTLMs      int64
}

func (b *Breaker) Trip() {
	b.State = OPEN
	b.Failures++
}

func (b *Breaker) Reset() {
	b.State = CLOSED
	b.Failures = 0
}
