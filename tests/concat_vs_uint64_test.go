package tests

import "testing"
import "strings"
import "strconv"

func BenchmarkUint64PK(b *testing.B) {
	m := make(map[uint64]struct{})

	for i := 0; i < b.N; i += 1 {
		var pk uint64 = uint64(i + i%10*10 + i%100*100 + i%1000*1000 + i%10000*10000 + i%100000*100000)
		m[pk] = struct{}{}
	}
}

func BenchmarkStringBuilderPK(b *testing.B) {
	m := make(map[string]struct{})

	for i := 0; i < b.N; i += 1 {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(strconv.Itoa(i * 10))
		sb.WriteString(strconv.Itoa(i * 100))
		sb.WriteString(strconv.Itoa(i * 1000))
		sb.WriteString(strconv.Itoa(i * 10000))
		sb.WriteString(strconv.Itoa(i * 100000))

		m[sb.String()] = struct{}{}
	}
}
