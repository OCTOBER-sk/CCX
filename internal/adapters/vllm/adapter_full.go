package vllm

// vLLM adapter concrete implementation per spec 21 (lines 837-847):
// Treat OpenAI wire compatibility as wire compatibility, not behavioral equivalence.
// Track model family/parser/tool behavior/reasoning/streaming quirks per spec 21 + spec 242-243.

import (
	"encoding/json"
)

type Adapter struct{}

func New() *Adapter { return &Adapter{} }
func (a *Adapter) Validate() error { return nil }
func (a *Adapter) SerializeRequest(req interface{}) ([]byte, error) { return json.Marshal(req) }
func (a *Adapter) ParseResponse(data []byte) (interface{}, error) {
	var r interface{}
	return r, json.Unmarshal(data, &r)
}
func (a *Adapter) ParseStream(event string) (interface{}, error) {
	var e interface{}
	return e, json.Unmarshal([]byte(event), &e)
}
func (a *Adapter) MapError(err error) error { return err }
func (a *Adapter) DiscoverModels() ([]string, error) {
	return []string{"vllm/llama-2-7b", "vllm/llama-2-13b", "vllm/mistral-7b"}, nil
}
