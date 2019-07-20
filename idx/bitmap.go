package idx

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/dicts"
	"sync"
)

type Bitmap struct {
	Mx sync.RWMutex
	B  *roaring.Bitmap
}

type HashBitmap struct {
	string  map[string]*Bitmap
	uint32  map[uint32]*Bitmap
	sex1    *Bitmap
	sex2    *Bitmap
	status1 *Bitmap
	status2 *Bitmap
	status3 *Bitmap
	mx      sync.RWMutex
}

func NewBitmap() *Bitmap {
	return &Bitmap{B: roaring.New()}
}

func NewBitmapFromSlice(s []uint32) *Bitmap {
	b := &Bitmap{B: roaring.New()}
	b.B.AddMany(s)
	return b
}

// Это нужно для начального несуществующего битмэпа, с которой будет идти объединение
func NewNullBitmap() *Bitmap {
	return &Bitmap{}
}

func (b *Bitmap) Add(x uint32) {
	b.Mx.Lock()
	b.B.Add(x)
	b.Mx.Unlock()
}

func (b *Bitmap) AddUnsafe(x uint32) {
	b.B.Add(x)
}

func (b *Bitmap) Remove(x uint32) {
	b.Mx.Lock()
	b.B.Remove(x)
	b.Mx.Unlock()
}

func (b *Bitmap) Cardinality() uint64 {
	b.Mx.Lock()
	c := b.B.GetCardinality()
	b.Mx.Unlock()
	return c
}

func (r *Bitmap) Clone() (*Bitmap, bool) {
	if r == nil {
		return nil, false
	}

	c := r.B.Clone()

	if c.IsEmpty() {
		return nil, false
	}

	return &Bitmap{
		B: c,
	}, true
}

func (r *Bitmap) And(and *Bitmap) (*Bitmap, bool) {
	// base bitmap
	if r == nil || r.B == nil {
		if and == nil {
			return nil, false
		} else {
			return and.Clone()
		}
	}

	if and == nil {
		return nil, false
	}

	r.B.And(and.B)

	if r.B.IsEmpty() {
		return nil, false
	}

	return r, true
}

func (r *Bitmap) Or(or *Bitmap) (*Bitmap, bool) {
	if r == nil || r.B == nil {
		if or == nil {
			return nil, false
		} else {
			return or.Clone()
		}
	}

	if or == nil {
		return r, true
	}

	r.B.Or(or.B)

	return r, true
}

func (r *Bitmap) AndNot(not *Bitmap) (*Bitmap, bool) {
	if not == nil {
		return r, true
	}

	r.B.AndNot(not.B)

	if r.B.IsEmpty() {
		return nil, false
	}

	return r, true
}

func (r *Bitmap) AndOr(and, or *Bitmap) (*Bitmap, bool) {
	// base bitmap always nil
	if r.B == nil {
		if and != nil {
			if s, ok := and.Clone(); ok {
				r.B = s.B

				return r.Or(or)
			} else {
				return nil, false
			}
		}
	}
	if and == nil {
		return r.And(or)
	}

	if s, ok := or.Clone(); ok {
		s, _ = s.Or(and)

		return r.And(s)
	}

	return r.And(and)
}

func NewStringHashBitmap() *HashBitmap {
	return &HashBitmap{string: make(map[string]*Bitmap)}
}

func NewUint32HashBitmap() *HashBitmap {
	return &HashBitmap{uint32: make(map[uint32]*Bitmap)}
}

func NewSexHashBitmap() *HashBitmap {
	return &HashBitmap{
		sex1: NewBitmap(),
		sex2: NewBitmap(),
	}
}

func NewStatusHashBitmap() *HashBitmap {
	return &HashBitmap{
		status1: NewBitmap(),
		status2: NewBitmap(),
		status3: NewBitmap(),
	}
}

func (d *HashBitmap) Uint32Cardinality() int {
	return len(d.uint32)
}

func (b *HashBitmap) Uint32PrintKeys() {
	b.mx.Lock()
	for k := range b.uint32 {
		println(k)
	}
	b.mx.Unlock()

}

func (d *HashBitmap) SexGet(sex dicts.Sex) *Bitmap {
	if sex == dicts.SexMale {
		return d.sex1
	}

	return d.sex2
}

func (d *HashBitmap) StatusGet(status dicts.Status) *Bitmap {
	if status == dicts.StatusFree {
		return d.status1
	} else if status == dicts.StatusOccupied {
		return d.status2
	}

	return d.status3
}

func (d *HashBitmap) StringGet(k string) *Bitmap {
	d.mx.RLock()
	v, ok := d.string[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	return nil
}

func (d *HashBitmap) StringGetOrCreate(k string) *Bitmap {
	d.mx.RLock()
	v, ok := d.string[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	r := &Bitmap{B: roaring.New()}
	d.mx.Lock()
	d.string[k] = r
	d.mx.Unlock()
	return r
}

func (d *HashBitmap) Uint32GetOrCreate(k uint32) *Bitmap {
	d.mx.RLock()
	v, ok := d.uint32[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	r := &Bitmap{B: roaring.New()}
	d.mx.Lock()
	d.uint32[k] = r
	d.mx.Unlock()
	return r
}

func (d *HashBitmap) Uint32Get(k uint32) *Bitmap {
	d.mx.RLock()
	v, ok := d.uint32[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	return nil

}
