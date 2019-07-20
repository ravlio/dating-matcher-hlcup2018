package skiplist

import (
	"github.com/ravlio/highloadcup2018/idx"
	"github.com/ravlio/highloadcup2018/skiplist/int32t"
	"github.com/ravlio/highloadcup2018/skiplist/stringt"
	"github.com/rs/zerolog/log"
	"sync"
)

type Skiplist struct {
	Mx     sync.RWMutex
	int32  *int32t.Skiplist
	string *stringt.Skiplist
}

func NewInt32Skiplist(cmp int32t.Comparator) *Skiplist {
	return &Skiplist{int32: int32t.New(cmp)}
}

func NewStringSkiplist(cmp stringt.Comparator) *Skiplist {
	return &Skiplist{string: stringt.New(cmp)}
}

func (s *Skiplist) Int32Insert(k int32, acc uint32) {
	s.Mx.Lock()
	_, err := s.int32.Insert(k, acc)
	s.Mx.Unlock()

	if err != nil {
		log.Error().Err(err)
	}
}

func (s *Skiplist) Int32Delete(k int32, acc uint32) bool {
	return s.int32.DeleteValue(k, acc)
}

func (s *Skiplist) Int32Select(k int32) (bm *idx.Bitmap, found bool, err error) {
	return s.int32.SelectFromBitmap(k)
}

func (s *Skiplist) StringInsert(k string, acc uint32) {
	s.Mx.Lock()
	_, err := s.string.Insert(k, acc)

	if err != nil {
		log.Error().Err(err)
	}

	s.Mx.Unlock()
}

func (s *Skiplist) StringDelete(k string, acc uint32) bool {
	return s.string.DeleteValue(k, acc)
}

func (s *Skiplist) StringSelect(k string) (bm *idx.Bitmap, found bool, err error) {
	return s.string.SelectFromBitmap(k)
}
