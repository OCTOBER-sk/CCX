package openai_responses

// OpenAI Responses adapter — separate adapter, NOT Chat++ per spec 113 (lines 3875-3893).
// Remains a distinct protocol module with its own serialize/parse logic.

import (
	"github.com/ccx/ccx/internal/protocol"
)

type Adapter struct{}

func New() protocol.Adapter {
	return &Adapter{}
}

func (a *Adapter) Validate() error { return nil }
func (a *Adapter) SerializeRequest(req interface{}) ([]byte, error) { return nil, nil }
func (a *Adapter) ParseResponse(data []byte) (interface{}, error) { return nil, nil }
func (a *Adapter) ParseStream(event string) (interface{}, error) { return nil, nil }
func (a *Adapter) MapError(err error) error { return err }
func (a *Adapter) DiscoverModels() ([]string, error) { return []string{}, nil }
