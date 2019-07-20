package tests

import (
	"fmt"
	"github.com/dghubble/trie"
	// "github.com/derekparker/trie"
	"github.com/pkg/errors"
	"github.com/ravlio/highloadcup2018/trie/patricia"
	"github.com/ravlio/highloadcup2018/utils"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestUnicodePatricia(t *testing.T) {
	tr := patricia.NewTrie()
	tr.Insert(patricia.Prefix("привет"), 2)
	tr.Insert(patricia.Prefix("привет2"), 1)
	tr.Insert(patricia.Prefix("прив"), 3)
	tr.Insert(patricia.Prefix("при"), 4)
	tr.Insert(patricia.Prefix("здорова"), 5)

	i := 0
	r := []uint32{4, 3, 2, 1}
	err := tr.VisitSubtree(patricia.Prefix("при"), func(prefix patricia.Prefix, item patricia.Item) error {
		if uint32(item) != r[i] {
			return errors.New("not found")
		}

		i++
		return nil
	})
	if err != nil {
		t.Fail()
	}
}

func TestPatriciaFa(t *testing.T) {
	tr := patricia.NewTrie()
	tr.Insert(patricia.Prefix("Фаушуко"), 1)

	i := 0
	r := []uint32{1}
	err := tr.VisitSubtree(patricia.Prefix("Фа"), func(prefix patricia.Prefix, item patricia.Item) error {
		if uint32(item) != r[i] {
			return errors.New("not found")
		}

		i++
		return nil
	})
	if err != nil {
		t.Fail()
	}

	if i != 1 {
		t.Fail()
	}
}

func TestPatriciaDelete(t *testing.T) {
	tr := patricia.NewTrie()
	tr.Insert(patricia.Prefix("привет"), 1)
	tr.Insert(patricia.Prefix("привет"), 2)
	tr.Insert(patricia.Prefix("здорова"), 3)
	tr.Insert(patricia.Prefix("прив"), 3)
	f := false
	err := tr.VisitSubtree(patricia.Prefix("приве"), func(prefix patricia.Prefix, item patricia.Item) error {
		if string(prefix) == "привет" {
			f = true
		} else {
			return errors.New("not found")
		}
		return nil
	})
	if err != nil {
		t.Fail()
	}
	if !f {
		t.Fail()
	}
}

func BenchmarkPatriciaInsert(b *testing.B) {
	tr := patricia.NewTrie()

	for i := 0; i <= b.N; i++ {
		tr.Insert(patricia.Prefix("sdfdsf"), patricia.Item(i))
	}
}

func BenchmarkDghubbleInsert(b *testing.B) {
	tr := trie.NewRuneTrie()

	for i := 0; i <= b.N; i++ {
		tr.Put("sdfdsf", i)
	}
}

func BenchmarkPatricia(b *testing.B) {
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	t := patricia.NewTrie()

	for i := 0; i < 1000000; i++ {
		t.Insert(patricia.Prefix(strings.ToLower(utils.RandStringRunes(rand.Intn(20-5)+5))), patricia.Item(i))
	}

	str := make([]string, 1000000)

	for i := 0; i < 1000000; i++ {
		str[i] = strings.ToLower(utils.RandStringRunes(rand.Intn(6-3) + 3))
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Println((m2.Alloc - m1.Alloc) / 1024 / 1024)

	b.ResetTimer()
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := str[i]

			t.VisitSubtree(patricia.Prefix(s), func(prefix patricia.Prefix, item patricia.Item) error {
				return nil
			})

			i++
			if i >= 1000000 {
				i = 0
			}
		}
	})
}

/*
func BenchmarkDerekparkerTrie(b *testing.B) {
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	t := trie.New()

	for i := 0; i < 1000000; i++ {
		t.Add(strings.ToLower(utils.RandStringRunes(rand.Intn(20-5)+5)), i)
	}

	str := make([]string, 1000000)

	for i := 0; i < 1000000; i++ {
		str[i] = strings.ToLower(utils.RandStringRunes(rand.Intn(6-3) + 3))
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Println((m2.Alloc - m1.Alloc) / 1024 / 1024)

	b.ResetTimer()
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := str[i]

			t.PrefixSearch(s)

			i++
			if i >= 1000000 {
				i = 0
			}
		}
	})
}*/
