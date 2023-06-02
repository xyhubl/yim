package bytes

import (
	"reflect"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	w := NewWriterSize(0)
	for i := 0; i < 30; i++ {
		w.Write(make([]byte, 1024, 1024))
		t.Log(w.Size(), w.n)
	}
}

func TestWriter(t *testing.T) {
	w := NewWriterSize(64)
	if w.Len() != 0 && w.Size() != 64 {
		t.FailNow()
	}
	b := []byte("hello")
	w.Write(b)
	if !reflect.DeepEqual(b, w.Buffer()) {
		t.FailNow()
	}
	w.Peek(len(b))
	w.Reset()
	for i := 0; i < 1024; i++ {
		w.Write(b)
	}
	w.Reset()
	if w.Len() != 0 {
		t.FailNow()
	}
}
