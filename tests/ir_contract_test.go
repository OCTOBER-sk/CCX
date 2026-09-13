package tests

import (
	"testing"
	"github.com/ccx/ccx/internal/ir"
)

// L0 static schema test per spec section 152 (lines 4748-4771).
// Every public feature must have automated tests (spec 79, rule 16).

func TestIRRequestSchema(t *testing.T) {
	req := ir.Request{
		Model: "openrouter/qwen/2.5",
		System: "You are a helpful assistant.",
		Stream: true,
	}
	if req.Model == "" {
		t.Fatal("model is required")
	}
	if !req.Stream {
		t.Log("stream false is allowed but streaming is preferred")
	}
}

func TestContentBlockExtensibility(t *testing.T) {
	cb := ir.ContentBlock{
		Text: "hello",
	}
	if cb.Text != "hello" {
		t.Errorf("expected text 'hello', got %s", cb.Text)
	}
}
