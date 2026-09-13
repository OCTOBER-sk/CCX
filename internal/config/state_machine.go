package config

// Config state machine per spec 101 (lines 3365-3430): READ → VALIDATE → PLAN → SNAPSHOT → APPLY → VERIFY → COMMIT.
// Failure at any stage → ROLLBACK. Transactional (no partial mutations by default per spec 107, 503-504).

type State int

const (
	READ State = iota
	VALIDATE
	PLAN
	SNAPSHOT
	APPLY
	VERIFY
	COMMIT
	ROLLBACK
)

type Transaction struct {
	CurrentState State
	Snapshot string
	Applied bool
	Verified bool
}

func (t *Transaction) Start() {
	t.CurrentState = READ
}

func (t *Transaction) Advance() bool {
	switch t.CurrentState {
	case READ:
		t.CurrentState = VALIDATE
	case VALIDATE:
		t.CurrentState = PLAN
	case PLAN:
		t.CurrentState = SNAPSHOT
	case SNAPSHOT:
		t.CurrentState = APPLY
	case APPLY:
		t.CurrentState = VERIFY
	case VERIFY:
		t.Applied = true
		t.CurrentState = COMMIT
	case COMMIT:
		t.Verified = true
		return false // completed
	case ROLLBACK:
		t.Applied = false
		t.Verified = false
		return false
	}
	return true
}

func (t *Transaction) Rollback() {
	t.CurrentState = ROLLBACK
	t.Applied = false
	t.Verified = false
}
