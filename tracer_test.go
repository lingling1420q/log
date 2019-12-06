package log

import (
	"testing"

	"log/trace"
)

func TestIsTracer(t *testing.T) {
	lg := New()
	_, ok := lg.(trace.Tracer)
	if !ok {
		t.Error("want true")
		return
	}
}
