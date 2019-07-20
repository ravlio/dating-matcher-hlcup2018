package tests

import (
	"github.com/RoaringBitmap/roaring"
	// "github.com/google/btree"
	"github.com/ravlio/highloadcup2018/utils"
	"math/rand"
	"reflect"
	"testing"
	"time"
)
import int32s "github.com/ravlio/highloadcup2018/skiplist/int32t"
import strs "github.com/ravlio/highloadcup2018/skiplist/stringt"
import fskiplist "github.com/sean-public/fast-skiplist"

// import mgsl "github.com/MauriceGit/skiplist"

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func init() {

}

func TestInt32SkiplistGt(t *testing.T) {
	s := int32s.New(int32s.BuiltinGreaterThan)
	s.Insert(1, 1)
	s.Insert(1, 2)
	s.Insert(1, 2)
	s.Insert(2, 3)
	s.Insert(2, 4)
	s.Insert(5, 5)
	s.Insert(6, 6)
	s.Insert(4, 7)

	i, _ := s.SelectFrom(1)

	type r struct {
		k int32
		v uint32
	}

	var rs = make([]r, 0, 7)

	for i.Next() {
		rs = append(rs, r{i.Key(), i.Value()})

	}

	exp := []r{
		//{1, 2},
		//{1, 2},
		//{1, 1},
		{2, 3},
		{2, 4},
		{4, 7},
		{5, 5},
		{6, 6},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.FailNow()
	}

	i, _ = s.SelectFrom(4)

	rs = make([]r, 0, 3)

	for i.Next() {
		rs = append(rs, r{i.Key(), i.Value()})
	}

	exp = []r{
		//{4, 7},
		{5, 5},
		{6, 6},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.Fail()
	}
}

func TestInt32SkiplistLt(t *testing.T) {
	s := int32s.New(int32s.BuiltinLessThan)
	s.Insert(1, 1)
	s.Insert(1, 2)
	s.Insert(1, 2)
	s.Insert(2, 3)
	s.Insert(2, 4)
	s.Insert(5, 5)
	s.Insert(6, 6)
	s.Insert(4, 7)

	i, _ := s.SelectFrom(6)

	type r struct {
		k int32
		v uint32
	}

	var rs = make([]r, 0, 7)

	for i.Next() {
		rs = append(rs, r{i.Key(), i.Value()})

	}

	exp := []r{
		//{6, 6},
		{5, 5},
		{4, 7},
		{2, 3},
		{2, 4},
		{1, 1},
		{1, 2},
		{1, 2},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.FailNow()
	}

	i, _ = s.SelectFrom(4)

	rs = make([]r, 0, 3)

	for i.Next() {
		rs = append(rs, r{i.Key(), i.Value()})

	}

	exp = []r{
		//{4, 7},
		{2, 3},
		{2, 4},
		{1, 1},
		{1, 2},
		{1, 2},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.Fail()
	}
}

func TestInt32SkiplistLessThanBitmap(t *testing.T) {
	s := int32s.New(int32s.BuiltinLessThan)
	s.Insert(11, 1)
	s.Insert(11, 2)
	s.Insert(11, 3)
	s.Insert(12, 4)
	s.Insert(12, 5)
	s.Insert(15, 6)
	s.Insert(16, 7)
	s.Insert(14, 8)

	i, _, _ := s.SelectFromBitmap(12)

	b := roaring.NewBitmap()
	b.AddMany([]uint32{1, 2, 3})

	println(i.B.String())
	if b.String() != i.B.String() {
		t.FailNow()
	}
}

func TestInt32SkiplistGreaterThanBitmap(t *testing.T) {
	s := int32s.New(int32s.BuiltinGreaterThan)
	s.Insert(11, 1)
	s.Insert(11, 2)
	s.Insert(11, 3)
	s.Insert(12, 4)
	s.Insert(12, 5)
	s.Insert(15, 6)
	s.Insert(16, 7)
	s.Insert(14, 8)

	i, _, _ := s.SelectFromBitmap(12)

	b := roaring.NewBitmap()
	b.AddMany([]uint32{6, 7, 8})

	if b.String() != i.B.String() {
		t.FailNow()
	}
}

/*func TestFastSkiplistLessThen(t *testing.T) {
	s := fskiplist.New()
	s.Set(1, 1)
	s.Set(1, 2)
	s.Set(1, 2)
	s.Set(2, 3)
	s.Set(2, 4)
	s.Set(5, 5)
	s.Set(6, 6)
	s.Set(4, 7)

	i := s.Get(4)

	type r struct {
		k float64
		v interface{}
	}

	var rs = make([]r, 0, 7)

	for i.Next() != nil {
		rs = append(rs, r{i.
			Key(), i.Value()})

	}

	exp := []r{
		//{1, 2},
		//{1, 2},
		//{1, 1},
		{2, 3},
		{2, 4},
		{4, 7},
		{5, 5},
		{6, 6},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.FailNow()
	}

	i = s.Get(4)

	rs = make([]r, 0, 3)

	for i.Next() != nil {
		rs = append(rs, r{i.Key(), i.Value()})
	}

	exp = []r{
		//{4, 7},
		{5, 5},
		{6, 6},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.Fail()
	}
}*/

func TestStringSkiplistGt(t *testing.T) {
	s := strs.New(strs.BuiltinGreaterThan)
	s.Insert("trulala@gmail.com", 1)
	s.Insert("1au@gmail.com", 2)
	s.Insert("1b@mail.ru", 3)
	s.Insert("2a@mail.ru", 4)
	s.Insert("5@mail.ru", 5)
	s.Insert("zzzz1@mail.ru", 6)
	s.Insert("22zz@gmail.com", 7)
	s.Insert("4@gmail.com", 8)

	i, _ := s.SelectFrom("1b@mail.r")

	type r struct {
		k string
		v uint32
	}

	var rs = make([]r, 0, 7)

	for i.Next() {
		rs = append(rs, r{i.Key(), i.Value()})

	}

	exp := []r{
		//{"1au@gmail.com", 2},
		{"1b@mail.ru", 3},
		{"22zz@gmail.com", 7},
		{"2a@mail.ru", 4},
		{"4@gmail.com", 8},
		{"5@mail.ru", 5},
		{"trulala@gmail.com", 1},
		{"zzzz1@mail.ru", 6},
	}

	if !reflect.DeepEqual(exp, rs) {
		t.FailNow()
	}
}

func TestStringSkiplistLtCase1(t *testing.T) {
	s := strs.New(strs.BuiltinLessThan)

	s.Insert("1sdfsdf", 2)
	s.Insert("zdsfsdfsdf", 3)
	s.Insert("zfdssdfsdf", 4)
	s.Insert("zfdfdfdf", 5)

	i, _, _ := s.SelectFromBitmap("et")

	b := roaring.NewBitmap()
	b.AddMany([]uint32{2})

	if b.String() != i.B.String() {
		t.FailNow()
	}

	s.Insert("odresiol@ya.ru", 1)
	s.Insert("aewerwer", 6)

	b = roaring.NewBitmap()
	b.AddMany([]uint32{2, 6})

	i, _, _ = s.SelectFromBitmap("et")

	if b.String() != i.B.String() {
		t.FailNow()
	}

}

func TestStringSkiplistDelete(t *testing.T) {
	s := strs.New(strs.BuiltinLessThan)

	s.Insert("1sdfsdf", 2)
	s.Insert("1sdfsdff", 6)
	s.Insert("zdsfsdfsdf", 3)
	s.Insert("zfdssdfsdf", 4)
	s.Insert("zfdfdfdf", 5)

	i, _, _ := s.SelectFromBitmap("et")

	s.PrintStats()
	b := roaring.NewBitmap()
	b.AddMany([]uint32{2, 6})

	if b.String() != i.B.String() {
		t.FailNow()
	}

	if s.DeleteValue("1sdfsdf", 3) != false {
		t.FailNow()
	}

	if s.DeleteValue("sdfsdf", 2) != false {
		t.FailNow()
	}

	if s.DeleteValue("1sdfsdf", 2) != true {
		t.FailNow()
	}

	i, _, _ = s.SelectFromBitmap("et")

	s.PrintStats()
	b = roaring.NewBitmap()
	b.AddMany([]uint32{6})

	println(i.B.String())
	if b.String() != i.B.String() {
		t.FailNow()
	}

}

/*func BenchmarkInt32SkiplistInsertTimeDescendingM(b *testing.B) {
	var i int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			if i >= 10000000 {
				i = 0
			}
			if _, err := sk.Insert(skKeys[i], i); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkInt32SkiplistInsertTimeDescendingS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := sk.Insert(skKeys[i], int32(i)); err != nil {
			b.Fatal(err)
		}
	}
}*/

func BenchmarkInt32SkiplistLess(b *testing.B) {
	sk := int32s.New(int32s.BuiltinLessThan)
	for i := 0; i < 1000000; i++ {
		sk.Insert(int32(random(0, 473385600)), uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := int32(random(1, 473385600))
		b.StartTimer()
		_, err := sk.SearchFrom(r, 10)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInt32SkiplistGreater(b *testing.B) {
	sk := int32s.New(int32s.BuiltinGreaterThan)
	for i := 0; i < 1000000; i++ {
		sk.Insert(int32(random(0, 473385600)), uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := int32(random(1, 473385600))
		b.StartTimer()
		_, err := sk.SearchFrom(r, 10)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFastSkiplistLess(b *testing.B) {
	list := fskiplist.New()

	for i := 0; i < 1000000; i++ {
		list.Set(float64(random(0, 473385600)), i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := float64(random(1, 473385600))
		b.StartTimer()
		e := list.Get(r)
		if e != nil {
			c := 0
			for e.Next() != nil {
				c++
				if c > 10 {
					break
				}
			}
		}
	}
}

func BenchmarkStringSkiplistLess(b *testing.B) {
	sk := strs.New(strs.BuiltinLessThan)
	for i := 0; i < 1000000; i++ {
		sk.Insert(utils.RandStringRunes(random(1, 100)), uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := sk.SelectFromBitmap(utils.RandStringRunes(random(1, 100)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringSkiplistGreater(b *testing.B) {
	sk := strs.New(strs.BuiltinGreaterThan)
	for i := 0; i < 100000; i++ {
		sk.Insert(utils.RandStringRunes(random(1, 100)), uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := sk.SelectFromBitmap(utils.RandStringRunes(random(1, 100)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringSkiplist1mInsert(b *testing.B) {
	sk := strs.New(strs.BuiltinGreaterThan)

	b.Run("1-10k", func(b *testing.B) {
		// time.Sleep(time.Second)
		for i := 0; i <= 10000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("10k-50k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 10001; i <= 50000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("50k-200k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 50001; i <= 200000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("200k-500k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 200001; i <= 500000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("500k-1m", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 500001; i <= 1000000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})
}

func BenchmarkInt32Skiplist1mInsert(b *testing.B) {
	var gi uint32
	sk := int32s.New(int32s.BuiltinGreaterThan)
	for j := 0; j < 5; j++ {
		b.Run("1", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i <= b.N; i++ {
				b.StopTimer()
				s := int32(random(0, 473385600))
				b.StartTimer()
				sk.Insert(s, gi)
				gi++
			}
		})
	}
}

func BenchmarkFastSkiplist1mInsert(b *testing.B) {
	var gi int
	list := fskiplist.New()
	for j := 0; j < 5; j++ {
		b.Run("1", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i <= b.N; i++ {
				b.StopTimer()
				s := int32(random(0, 473385600))
				b.StartTimer()
				list.Set(float64(s), gi)
				gi++
			}
		})
	}
}

func BenchmarkBtree1mInsert(b *testing.B) {
	sk := strs.New(strs.BuiltinGreaterThan)

	b.Run("1-10k", func(b *testing.B) {
		// time.Sleep(time.Second)
		for i := 0; i <= 10000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("10k-50k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 10001; i <= 50000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("50k-200k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 50001; i <= 200000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("200k-500k", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 200001; i <= 500000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})

	b.Run("500k-1m", func(b *testing.B) {
		time.Sleep(time.Second)
		for i := 500001; i <= 1000000; i++ {
			b.StopTimer()
			s := utils.RandStringRunes(random(1, 100))
			b.StartTimer()
			sk.Insert(s, uint32(i))
		}
	})
}

/*
func BenchmarkBtreeGreater(b *testing.B) {
	bt := btree.New(32)
	for i := 0; i < 1000000; i++ {

		bt.ReplaceOrInsert(btree.Int(random(0, 473385600)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := btree.Int(random(0, 473385600))
		b.StartTimer()
		bt.DescendGreaterThan(r, func(i btree.Item) bool {
			return true
		})
	}
}

func BenchmarkBtreeLess(b *testing.B) {
	bt := btree.New(32)
	for i := 0; i < 100000; i++ {

		bt.ReplaceOrInsert(btree.Int(random(0, 473385600)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := btree.Int(random(0, 473385600))
		b.StartTimer()
		bt.AscendLessThan(r, func(i btree.Item) bool {
			return true
		})
	}
}*/
