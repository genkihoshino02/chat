package trace

import (
	"testing"
	"bytes"
)

func TestNew(t *testing.T){
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil{
		t.Error("return value from tracer is nil ")
	}else{
		tracer.Trace("Hello, trace package")
		if buf.String() != "Hello, trace package\n"{
			t.Errorf("'%s' is wrong output",buf.String())
		}
	}
}

func TestOff(t *testing.T){
	var silentTracer Tracer = Off()
	silentTracer.Trace("Data")
}