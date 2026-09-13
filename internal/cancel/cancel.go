package cancel

// Cancellation propagation per spec section 30, lines 1147-1161.
// Propagation chain: Claude Code (line 1154) -> CCX cancellation (line 1156) -> provider cancellation (line 1158).
// Per line 1161: CCX MUST not leave abandoned generations running unnecessarily.

// Propagation tracks cancellation through the three-layer chain.
type Propagation struct {
	Active bool // line 1156: CCX cancellation active
	UpstreamCanceled bool // line 1158: provider cancellation signaled
}

// Cancel initiates the propagation chain per lines 1154-1158.
func (p *Propagation) Cancel() {
	p.Active = true
	p.UpstreamCanceled = true
}

// IsCanceled reports whether cancellation is in progress at any layer.
func (p *Propagation) IsCanceled() bool {
	return p.Active || p.UpstreamCanceled
}
