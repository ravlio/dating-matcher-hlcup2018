package tests

import "testing"

type D interface {
	Append(D)
}

type Strings []string

func (s Strings) Append(d D) {}

func BenchmarkInterface(b *testing.B) {
	s := D(Strings{})
	for i := 0; i < b.N; i += 1 {
		s.Append(Strings{""})
	}
}

func BenchmarkConcrete(b *testing.B) {
	s := Strings{} // only difference is that I'm not casting it to the generic interface
	for i := 0; i < b.N; i += 1 {
		s.Append(Strings{""})
	}
}
