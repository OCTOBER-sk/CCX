package transform

// Transform engine per spec 4.9 (lines 247-258) and 116 (lines 3956-3968):
// Policy operations: PASS, CONVERT, REMOVE, EMULATE, DEGRADE, REJECT.
// Every transform must declare: id, input, output, loss, conditions, reversible, risk, default.
// Never hide semantic loss.

type Policy int

const (
	PASS Policy = iota
	CONVERT
	REMOVE
	EMULATE
	DEGRADE
	REJECT
)

type Transform struct {
	ID          string
	Input       string
	Output      string
	Loss        string
	Reversible  bool
	Risk        string
	Default     Policy
}
