package llamacpp

// Adapter concrete implementation per spec 23/21 (provider adapter framework).
import ("encoding/json")

type Adapter struct{}
func New() *Adapter { return &Adapter{} }
func (a *Adapter) Validate() error { return nil }
func (a *Adapter) SerializeRequest(req interface{}) ([]byte, error) { return json.Marshal(req) }
func (a *Adapter) ParseResponse(data []byte) (interface{}, error) { var r interface{}; return r, json.Unmarshal(data, &r) }
func (a *Adapter) ParseStream(event string) (interface{}, error) { var e interface{}; return e, json.Unmarshal([]byte(event), &e) }
func (a *Adapter) MapError(err error) error { return err }
func (a *Adapter) DiscoverModels() ([]string, error) { return []string{"llama-2-7b,llama-2-13b,llama-3-8b"}, nil }
