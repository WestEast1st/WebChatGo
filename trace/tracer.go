package trace

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事を記憶することのできるオブジェクトです。
type Tracer interface {
	Trace(...interface{})
}

type nilTracer struct {
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返します
func Off() Tracer {
	return &nilTracer{}
}

// tracer構造体を返還する
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
