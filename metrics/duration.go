package metrics

import "time"
import "sync/atomic"
import "fmt"

type Duration struct {
	name  string
	count uint32
	time  uint64
}

func NewDurarion(name string) *Duration {
	d := &Duration{name: name}
	go d.run()
	return d
}

func (d *Duration) run() {
	for {
		time.Sleep(time.Second)
		c := atomic.LoadUint32(&d.count)
		t := atomic.LoadUint64(&d.time)
		if c > 0 || t > 0 {
			atomic.StoreUint32(&d.count, 0)
			atomic.StoreUint64(&d.time, 0)

			fmt.Printf(
				"Duration %q: count: %d, avg: %.2fms\n",
				d.name,
				c,
				float64(t/uint64(c))/1000/1000,
			)
		}

	}
}

func (d *Duration) Write(dd time.Time) {
	atomic.AddUint32(&d.count, 1)
	atomic.AddUint64(&d.time, uint64(time.Since(dd)))
}
