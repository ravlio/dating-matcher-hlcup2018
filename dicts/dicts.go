package dicts

import (
	"sync"
	"sync/atomic"
)

type Uint8Dict struct {
	k  map[string]uint8
	v  []string
	mx sync.RWMutex
	c  uint32
}

func NewUint8Dict() *Uint8Dict {
	return &Uint8Dict{k: make(map[string]uint8), v: make([]string, 0, 255)}
}

func (d *Uint8Dict) GetOrCreateValue(k string) uint8 {
	d.mx.RLock()
	v, ok := d.k[k]
	d.mx.RUnlock()

	if ok {
		return v
	}

	d.mx.Lock()
	v2, ok := d.k[k]
	if ok {
		d.mx.Unlock()
		return v2
	}

	i := uint8(atomic.AddUint32(&d.c, 1))
	d.k[k] = i
	d.v = append(d.v, k)
	d.mx.Unlock()
	return i
}

func (d *Uint8Dict) GetKey(v uint8) string {
	d.mx.RLock()
	r := d.v[v-1]
	d.mx.RUnlock()
	return r
}

type Uint16Dict struct {
	k  map[string]uint16
	v  []string
	mx sync.RWMutex
	c  uint32
}

func NewUint16() *Uint16Dict {
	return &Uint16Dict{k: make(map[string]uint16), v: make([]string, 0, 1023)}
}

func (d *Uint16Dict) GetOrCreateValue(k string) uint16 {
	d.mx.RLock()
	v, ok := d.k[k]
	d.mx.RUnlock()
	if ok {
		return v
	}

	d.mx.Lock()
	v2, ok := d.k[k]
	if ok {
		d.mx.Unlock()
		return v2
	}

	i := uint16(atomic.AddUint32(&d.c, 1))
	d.k[k] = i
	d.v = append(d.v, k)
	d.mx.Unlock()
	return i
}

func (d *Uint16Dict) GetKey(v uint16) string {
	d.mx.RLock()
	r := d.v[v-1]
	d.mx.RUnlock()
	return r
}

type Uint32Dict struct {
	k  map[string]uint32
	v  []string
	Mx sync.RWMutex
	c  uint32
}

func NewUint32Dict() *Uint32Dict {
	return &Uint32Dict{k: make(map[string]uint32), v: make([]string, 0, 1023)}
}

func (d *Uint32Dict) GetValueUnsafe(k string) uint32 {
	v, ok := d.k[k]
	if ok {
		return v
	}

	v2, ok := d.k[k]
	if ok {
		return v2
	}

	i := atomic.AddUint32(&d.c, 1)
	d.k[k] = i
	d.v = append(d.v, k)

	return i
}

func (d *Uint32Dict) GetOrCreateValue(k string) uint32 {
	d.Mx.RLock()
	v, ok := d.k[k]
	d.Mx.RUnlock()
	if ok {
		return v
	}

	d.Mx.Lock()
	v2, ok := d.k[k]
	if ok {
		d.Mx.Unlock()
		return v2
	}

	i := atomic.AddUint32(&d.c, 1)
	d.k[k] = i
	d.v = append(d.v, k)
	d.Mx.Unlock()
	return i
}

func (d *Uint32Dict) GetValue(k string) (uint32, bool) {
	d.Mx.RLock()
	v, ok := d.k[k]
	d.Mx.RUnlock()
	if ok {
		return v, true
	}

	return 0, false
}

func (d *Uint32Dict) GetKeyUnsafe(v uint32) string {
	r := d.v[v-1]
	return r
}

func (d *Uint32Dict) GetKey(v uint32) string {
	d.Mx.RLock()
	r := d.v[v-1]
	d.Mx.RUnlock()
	return r
}
