package slist

import "testing"
import "github.com/stretchr/testify/assert"
import "math/rand"
import "time"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestInsert(t *testing.T) {
	slist := []uint32{}

	Insert(&slist, 1)
	Insert(&slist, 2)
	Insert(&slist, 2)
	Insert(&slist, 10)
	Insert(&slist, 100)
	Insert(&slist, 100)
	Insert(&slist, 99)
	Insert(&slist, 3)

	if !assert.Equal(t, []uint32{1, 2, 3, 10, 99, 100}, slist) {
		t.FailNow()
	}
}

func TestBug1(t *testing.T) {
	s1 := []uint32{1546, 2414, 2934, 4380, 5554, 7436, 7658, 7932, 8758, 8816, 11340, 11964, 15550, 15860, 16060, 16620, 16710, 16856, 18094, 18544, 19796, 20066, 20348, 20968, 21352, 21818, 22408, 22420, 23230, 25126, 25468, 25886, 27778, 28286, 28496, 29358, 29928, 30914, 31210}
	s2 := []uint32{72, 1016, 1886, 3336, 5428, 6160, 6738, 7210, 7528, 7592, 8578, 9020, 9456, 10450, 10512, 10768, 12366, 12530, 13678, 13830, 14590, 15250, 15644, 16060, 16536, 16926, 18256, 18262, 18502, 18634, 19386, 19866, 22438, 22672, 22950, 24758, 24782, 25398, 25672, 27908, 30014, 30348, 31334}
	r := And(s1, s2)
	if !assert.Equal(t, []uint32{16060}, r) {
		t.FailNow()
	}
}
func BenchmarkInsert(b *testing.B) {
	r := make([]uint32, 100000, 100000)
	for k := range r {
		r[k] = uint32(rand.Intn(100000) + 1)
	}
	slist := []uint32{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Insert(&slist, uint32(r[i%100000]))
	}
}

func TestAnd(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		a := []uint32{1, 4}
		b := []uint32{1, 2, 3, 4, 5}

		if !assert.Equal(t, []uint32{1, 4}, And(a, b)) {
			t.FailNow()
		}
	})

	t.Run("1", func(t *testing.T) {
		a := []uint32{1, 2, 3, 4, 5}
		b := []uint32{1, 3, 4, 5, 8}

		if !assert.Equal(t, []uint32{1, 3, 4, 5}, And(a, b)) {
			t.FailNow()
		}
	})

	t.Run("2", func(t *testing.T) {
		a := []uint32{}
		b := []uint32{1, 3, 4, 5, 8}

		if !assert.Nil(t, And(a, b)) {
			t.FailNow()
		}
	})

	t.Run("2", func(t *testing.T) {
		var a []uint32
		b := []uint32{}

		if !assert.Nil(t, And(a, b)) {
			t.FailNow()
		}
	})

	t.Run("2", func(t *testing.T) {
		a := []uint32{}
		b := []uint32{}

		if !assert.Nil(t, And(a, b)) {
			t.FailNow()
		}
	})

}

func BenchmarkAnd(b *testing.B) {
	m := make(map[int][]uint32, 1000)
	r1 := make([]int, 1000, 1000)
	r2 := make([]int, 1000, 1000)
	for k := range r1 {
		r1[k] = (rand.Intn(1000))
		r2[k] = (rand.Intn(1000))
	}

	for i := 0; i < 1000; i++ {
		var p []uint32
		for j := 0; j < 100; j++ {
			Insert(&p, uint32(rand.Intn(100)+1))
		}

		m[i] = p
	}

	println("ok")
	b.ReportAllocs()
	b.SetBytes(6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		And(m[r1[i%1000]], m[r2[i%1000]])
	}
}
