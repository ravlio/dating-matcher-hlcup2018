package tests

import "testing"

func X(p *int) {
	*p += 1
}

var N = 1000000000

func BenchmarkClosure(b *testing.B) {
	for n := 0; n < b.N; n++ {
		p := 0
		x := func() {
			p += 1
		}
		x()
	}
}

func BenchmarkAnonymous(b *testing.B) {
	for n := 0; n < b.N; n++ {
		p := 0
		func() {
			p += 1
		}()
	}
}

func BenchmarkFunction(b *testing.B) {
	for n := 0; n < b.N; n++ {
		p := 0
		X(&p)
	}
}
