package idx

import "sync"

type Hash struct {
	string    map[string]uint32
	uint32    map[uint32]uint32
	uint32Inv map[uint32]uint32

	int64    map[int64]uint32
	int64Inv map[uint32]int64
	Mx       sync.RWMutex
}

func NewStringHash() *Hash {
	return &Hash{string: make(map[string]uint32)}
}

func NewUint32Hash() *Hash {
	return &Hash{uint32: make(map[uint32]uint32)}
}

func NewInt64Hash() *Hash {
	return &Hash{int64: make(map[int64]uint32)}
}

func (s *Hash) StringGet(k string) (uint32, bool) {
	s.Mx.RLock()
	v, ok := s.string[k]

	s.Mx.RUnlock()

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) StringSet(k string, uid uint32) {
	s.Mx.Lock()
	s.string[k] = uid
	s.Mx.Unlock()
}

func (s *Hash) StringGetUnsafe(k string) (uint32, bool) {
	v, ok := s.string[k]

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) StringSetUnsafe(k string, uid uint32) {
	s.string[k] = uid
}

func (s *Hash) Uint32Get(k uint32) (uint32, bool) {
	s.Mx.RLock()
	v, ok := s.uint32[k]

	s.Mx.RUnlock()

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) Uint32Set(k uint32, uid uint32) {
	s.Mx.Lock()
	s.uint32[k] = uid
	s.Mx.Unlock()
}

func (s *Hash) Uint32Cardinality() int {
	return len(s.uint32)
}

func (s *Hash) Uint32GetUnsafe(k uint32) (uint32, bool) {
	v, ok := s.uint32[k]

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) Uint32Delete(k uint32) {
	s.Mx.Lock()
	delete(s.uint32, k)
	s.Mx.Unlock()
}

func (s *Hash) Uint32SetUnsafe(k uint32, uid uint32) {
	s.uint32[k] = uid
}

func (s *Hash) Uint32DeleteUnsafe(k uint32) {
	delete(s.uint32, k)
}

func (s *Hash) Uint32Change(new uint32, id uint32) {
	s.Mx.Lock()
	old, ok := s.uint32Inv[id]
	if ok {
		delete(s.uint32, old)
	}

	s.uint32[new] = id
	s.uint32Inv[id] = new
	s.Mx.Unlock()

}

////// Int64

func (s *Hash) Int64Get(k int64) (uint32, bool) {
	s.Mx.RLock()
	v, ok := s.int64[k]

	s.Mx.RUnlock()

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) Int64Cardinality() int {
	return len(s.int64)
}

func (s *Hash) Int64Set(k int64, uid uint32) {
	s.Mx.Lock()
	s.int64[k] = uid
	s.Mx.Unlock()
}

func (s *Hash) Int64GetUnsafe(k int64) (uint32, bool) {
	v, ok := s.int64[k]

	if !ok {
		return 0, false
	}

	return v, true
}

func (s *Hash) Int64Delete(k int64) {
	s.Mx.Lock()
	delete(s.int64, k)
	s.Mx.Unlock()
}

func (s *Hash) Int64SetUnsafe(k int64, uid uint32) {
	s.int64[k] = uid
}

func (s *Hash) Int64DeleteUnsafe(k int64) {
	delete(s.int64, k)
}

func (s *Hash) Int64Change(new int64, id uint32) {
	s.Mx.Lock()
	old, ok := s.int64Inv[id]
	if ok {
		delete(s.int64, old)
	}

	s.int64[new] = id
	s.int64Inv[id] = new
	s.Mx.Unlock()

}
