package dicts

import (
	"github.com/tchap/go-patricia/patricia"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func BenchmarkTrie(b *testing.B) {
	t := patricia.NewTrie()
	t.Insert(patricia.Prefix("sd"))
}

func BenchmarkDict(b *testing.B) {
	s := make([]string, 255)
	r := make([]uint8, 1000000)

	for i := 0; i < 1000000; i++ {
		r[i] = uint8(rand.Intn(255))
	}
	for i := 0; i < 255; i++ {
		s[i] = RandStringRunes(rand.Intn(100-3) + 3)
	}

	d := NewUint8Dict()
	j := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			j++
			if j >= 1000000 {
				j = 0
			}
		}
	})

}
