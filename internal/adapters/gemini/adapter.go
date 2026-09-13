package gemini

// Gemini adapter per spec 114 (lines 3904-3923): dedicated semantic mapper.
// NOT "OpenAI with renamed fields". Explicit support for function calling, tool results,
// system instruction, multimodal parts, thinking, streaming, finish reasons, usage.
// Unsupported semantics must be surfaced explicitly.

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
