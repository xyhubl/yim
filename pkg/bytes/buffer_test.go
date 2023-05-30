package bytes

import "testing"

func TestBuffer(t *testing.T) {
	p := NewPool(10, 100)
	p.Get()
}

func BenchmarkNewPool(b *testing.B) {
	p := NewPool(10, 100)
	for i := 0; i < b.N; i++ {
		bf := p.Get()
		p.Put(bf)
	}
}
