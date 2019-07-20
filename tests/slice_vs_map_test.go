package tests

import "testing"

type Dummy struct {
	ID int
}

func BenchmarkMapByKey(b *testing.B) {
	m := make(map[int]*Dummy)

	for i := 0; i < 1000000; i++ {
		m[i] = &Dummy{ID: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a, _ := m[i%1000000]
		a.ID = i + 1
	}
}

func BenchmarkSliceByKey(b *testing.B) {
	s := make([]*Dummy, 1000001)

	for i := 0; i < 1000000; i++ {
		s[i] = &Dummy{ID: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := s[i%1000000]
		a.ID = i + 1
	}
}
