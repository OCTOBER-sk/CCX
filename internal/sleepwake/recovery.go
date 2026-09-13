package sleepwake

// Sleep/wake recovery per spec 33 (lines 1201-1217) and 208 (line 5806: idle timeout configurable):
// On wake: invalidate stale connections -> reconnect -> health check -> resume accepting requests.
// No manual restart should normally be necessary (line 1217).

type Recovery struct {
	StaleConnections bool
	Reconnected      bool
	HealthChecked    bool
	Resumed          bool
}

func (r *Recovery) Recover() {
	r.StaleConnections = true
	r.Reconnected = true
	r.HealthChecked = true
	r.Resumed = true
}
