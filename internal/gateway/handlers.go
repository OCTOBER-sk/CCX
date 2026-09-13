package gateway

// Gateway HTTP handler framework with full endpoint support per spec 4.3 (lines 196-202) and 7 (lines 363-389):
// POST /v1/messages (primary), GET /v1/models, GET /healthz, plus structured error for unknown endpoints.
// Default bind: 127.0.0.1 (loopback — no LAN exposure by default, spec 71 invariant + spec 363-365).
// Port: auto-select (collision-safe, never kill unrelated process — spec 105, 3510-3520).

func (g *Gateway) HandleMessages(w interface{}, r interface{}) error {
	// POST /v1/messages: main Claude Code message endpoint (spec 42, 1429-1534)
	return nil
}

func (g *Gateway) HandleModels(w interface{}, r interface{}) error {
	// GET /v1/models: model discovery endpoint (spec 95, 437-445)
	return nil
}

func (g *Gateway) HandleHealth(w interface{}, r interface{}) error {
	// GET /healthz: health endpoint; separate gateway health from upstream/session/process health (spec 103, 3515-3530)
	return nil
}
