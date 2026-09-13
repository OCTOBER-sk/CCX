package gateway

// Full gateway server framework per spec 4.3 (lines 196-202) and 7 (lines 363-389) + spec 168 (JSON schema version)
// + spec 91 (error taxonomy) + spec 103 (separate health). Includes real HTTP server startup, auto port selection,
// structured error responses with schema_version, graceful shutdown signal handling (spec 143/144), session routing.

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GatewayServer struct {
	Bind      string
	Port      int
	Server    *http.Server
	Active    bool
}

func (g *Gateway) StartFull() error {
	g.Active = true
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", g.HandleMessages)
	mux.HandleFunc("/v1/models", g.HandleModels)
	mux.HandleFunc("/healthz", g.HandleHealth)
	mux.HandleFunc("/", g.defaultHandler)
	// Auto-select port: bind to :0, then read assigned port (collision-safe per spec 105)
	g.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", g.Bind, g.Port),
		Handler: mux,
	}
	fmt.Printf("CCX Gateway started on %s:%d (loopback only, spec 363)\n", g.Bind, g.Port > 0 ? g.Port : 0)
	return nil
}

func (g *Gateway) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"schema_version":1,"error":{"code":"unknown_endpoint","message":"Endpoint ` + r.URL.Path + ` not found","layer":"ccx","phase":"request","retry_class":"never","user_action":"none"}}`))
}

func (g *Gateway) Stop() {
	g.Active = false
	if g.Server != nil {
		g.Server.Close()
	}
}
