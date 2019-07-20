package idx

/*
import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/slist"
	"sync"
)

type Slist struct {
	Mx sync.RWMutex
	S []uint32
}

type HashSlist struct {
	uint32  map[uint32]Slist
	mx      sync.RWMutex
}

func NewSlist() *Slist {
	return &Slist{S: make([]uint32,0)}
}

func (s *Slist) Add(x uint32) {
	s.Mx.Lock()
	slist.Insert(&s.S,x)
	b.Mx.Unlock()
}

func (s *Slist) And(and *Slist) []uint32 {
	return slist.And(s.S, and.S)
}


func NewUint32HashSlist() *HashSlist {
	return &HashSlist{uint32: make(map[uint32]*Slist)}
}

func (d *HashSlist) Uint32GetOrCreate(k uint32) *Slist {
	d.mx.RLock()
	v, ok := d.uint32[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	r := NewSlist()
	d.mx.Lock()
	d.uint32[k] = r
	d.mx.Unlock()
	return r
}

func (d *HashSlist) Uint32Get(k uint32) *Slist {
	d.mx.RLock()
	v, ok := d.uint32[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	return nil

}*/
