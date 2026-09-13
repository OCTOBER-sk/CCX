package openai_chat

// OpenAI Chat adapter — separate from OpenAI Responses per spec 113 (lines 3875-3893).
// Must not be treated as Chat++. Shares reusable primitives but remains distinct.

import (
	"github.com/ccx/ccx/internal/protocol"
)

type Adapter struct{}

func New() protocol.Adapter {
	return &Adapter{}
}

func (a *Adapter) Validate() error {
	return nil
}
func (a *Adapter) SerializeRequest(req interface{}) ([]byte, error) {
	return nil, nil
}
func (a *Adapter) ParseResponse(data []byte) (interface{}, error) {
	return nil, nil
}
func (a *Adapter) ParseStream(event string) (interface{}, error) {
	return nil, nil
}
func (a *Adapter) MapError(err error) error {
	return err
}
func (a *Adapter) DiscoverModels() ([]string, error) {
	return []string{}, nil
}
