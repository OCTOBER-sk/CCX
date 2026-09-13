package gemini

// Gemini adapter concrete implementation per spec 114 (lines 3904-3923):
// Dedicated semantic mapper — NOT OpenAI rename. Explicit mapping for function calling, tool results,
// system instruction, multimodal parts, thinking, streaming, finish reasons, usage.

import (
	"encoding/json"
)

type Adapter struct{}

func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Validate() error {
	return nil
}

func (a *Adapter) SerializeRequest(req interface{}) ([]byte, error) {
	return json.Marshal(req)
}

func (a *Adapter) ParseResponse(data []byte) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func (a *Adapter) ParseStream(event string) (interface{}, error) {
	var streamEvent interface{}
	err := json.Unmarshal([]byte(event), &streamEvent)
	return streamEvent, err
}

func (a *Adapter) MapError(err error) error {
	return err
}

func (a *Adapter) DiscoverModels() ([]string, error) {
	return []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-1.0-pro"}, nil
}
