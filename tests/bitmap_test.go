package tests

import (
	"github.com/RoaringBitmap/roaring"
	pilosa "github.com/pilosa/pilosa/roaring"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Test100kX100kMapSlice(t *testing.T) {
	var bm = make(map[uint32][]uint32)

	t.Run("generate 100k bitmaps", func(b *testing.T) {
		for i := 1; i < 100000; i += 1 {
			bm[uint32(i)] = make([]uint32, 0, 100000/100)

			for j := 1; j < 100000; j += 100 {
				bm[uint32(i)] = append(bm[uint32(i)], uint32(j))
			}
		}
	})
}

func Test100kX100kRoaringBitmap(t *testing.T) {
	var bm = make(map[int]*roaring.Bitmap)

	t.Run("generate 100k bitmaps", func(b *testing.T) {
		for i := 1; i < 100000; i += 1 {
			bm[i] = roaring.New()
			for j := 1; j < 100000; j += 100 {
				bm[i].Add(uint32(j))
			}
		}
	})
}

func Test100kX100kPilosa(t *testing.T) {

	var bm = make(map[int]*pilosa.Bitmap)

	t.Run("generate 100k bitmaps", func(b *testing.T) {
		for i := 1; i < 100000; i += 1 {
			bm[i] = pilosa.NewBitmap()

			for j := 1; j < 100000; j += 100 {
				bm[i].Add(uint64(j))
			}
		}
	})
}

func BenchmarkBitmap(b *testing.B) {
	m := make(map[int]*roaring.Bitmap, 1000)
	r1 := make([]int, 1000, 1000)
	r2 := make([]int, 1000, 1000)
	for k := range r1 {
		r1[k] = (rand.Intn(1000))
		r2[k] = (rand.Intn(1000))
	}

	for i := 0; i < 1000; i++ {
		m[i] = roaring.New()
		for j := 0; j < 50000; j++ {
			m[i].Add(uint32(rand.Intn(50000) + 1))
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := m[r1[i%1000]].Clone()
		m[r2[i%1000]].And(c)
	}
}
