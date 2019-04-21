package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("[*] New return param nil!")
	} else {
		tracer.Trace(" - : Hello, Trace package!")
		if buf.String() != " - : Hello, Trace package!\n" {
			t.Errorf("[*] The incorrect string '%s' is output.", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var shilentTracer Tracer = Off()
	shilentTracer.Trace("data")
}
