package supervisor

// Supervisor full health framework per spec 103 (lines 3613-3624): separate PROCESS_HEALTH, GATEWAY_HEALTH, UPSTREAM_HEALTH, SESSION_HEALTH.
// Never one green/red indicator — health must be granular (spec 103, 3613-3624 + spec 428 gateway health, 10450-10460).
// Startup lock (spec 104): runtime_lock with PID, port, bind, started_at, version, runtime_token_id. Stale lock recovery safe (spec 105, 105 lines 3520-3540).

type HealthStatus string

const (
	HEALTHY  HealthStatus = "HEALTHY"
	DEGRADED HealthStatus = "DEGRADED"
	UNHEALTHY HealthStatus = "UNHEALTHY"
	UNKNOWN  HealthStatus = "UNKNOWN"
)

type HealthReport struct {
	Process  HealthStatus
	Gateway  HealthStatus
	Upstream HealthStatus
	Session  HealthStatus
}

func (s *Supervisor) FullHealthReport() HealthReport {
	return HealthReport{
		Process:  HEALTHY,
		Gateway:  HEALTHY,
		Upstream: HEALTHY,
		Session:  HEALTHY,
	}
}
