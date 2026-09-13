package supervisor

import (
	"fmt"
)

// Supervisor per spec section 4.2 (lines 185-194): runtime lifecycle, process/service state, health, startup, shutdown, crash recovery.
// Health states: PROCESS_HEALTH, GATEWAY_HEALTH, UPSTREAM_HEALTH, SESSION_HEALTH (section 103, lines 3613-3624).
// Must handle startup locking (section 104), stale lock recovery (section 105), and crash recovery (section 32).

type Supervisor struct {
	Running bool
	LockFile string
	PID int
}

func New() *Supervisor {
	return &Supervisor{Running: false, LockFile: "runtime/ccx.lock", PID: 0}
}

func (s *Supervisor) Start() error {
	if s.Running {
		return fmt.Errorf("already running")
	}
	s.Running = true
	s.PID = 1 // stub
	fmt.Println("CCX Supervisor started")
	return nil
}

func (s *Supervisor) RecoverStaleLock() error {
	// Per spec 104 (lines 3410-3430) + spec 105 (port ownership / stale lock recovery): if lock file exists but process not running, recover by removing stale lock.
	return fmt.Errorf("stale lock recovery: would recover %s if process %d not running (spec 104/105)", s.LockFile, s.PID)
}

func (s *Supervisor) Health() map[string]string {
	// Per spec 103 (lines 3613-3624) + spec 4.2: health separate for process, gateway, upstream, session.
	health := map[string]string{
		"PROCESS": "HEALTHY",
		"GATEWAY": "HEALTHY",
		"UPSTREAM": "HEALTHY",
		"SESSION": "HEALTHY",
	}
	if !s.Running {
		health["PROCESS"] = "DOWN"
		health["SESSION"] = "DOWN"
	}
	if s.LockFile != "" && s.PID == 0 {
		health["PROCESS"] = "LOCK_STALE"
	}
	return health
}

func (s *Supervisor) Stop() error {
	s.Running = false
	fmt.Println("CCX Supervisor stopped")
	return nil
}
