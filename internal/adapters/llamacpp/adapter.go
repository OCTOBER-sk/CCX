package llamacpp

// llama.cpp adapter per spec 21 (lines 849-851): use actual exposed API and capability metadata.

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
