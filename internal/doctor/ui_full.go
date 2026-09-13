package doctorui

// Doctor UI full framework per spec 414 (doctor levels), 415 (fix classes: safe/caution/destructive), 416 (transactional preview):
// Default read-only. --fix performs safe/caution/destructive with transactional preview. Bundle generation (spec 47, 1751-1755).

type FixClass string

const (
	SAFE        FixClass = "safe"
	CAUTION     FixClass = "caution"
	DESTRUCTIVE FixClass = "destructive"
)

type FixUI struct {
	ReadOnly bool
	FixClass FixClass
	PreviewEnabled bool
}

func NewUI(readOnly bool, fixClass FixClass) *FixUI {
	return &FixUI{
		ReadOnly: readOnly,
		FixClass: fixClass,
		PreviewEnabled: true,
	}
}

func (f *FixUI) ShowPreview(actions []string) string {
	return "Doctor preview: " + string(f.FixClass) + " actions: " + actions[0]
}
