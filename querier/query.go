package querier

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/idx"
)

// Индексатор и процессор условий добавления/удаления индексов
// По хорошему, чтобы сэкономить на колбэках,
// условия лучше заинлайнить, но руками очень сложно не ошибиться.
// Возможно тут помог бы генератор кода
// Плюс такой подход позволяет паралелить индексацию

type Job struct {
	Name   string
	Bitmap func() *idx.Bitmap
	Hash   func() (aid uint32, found bool)
}

type Jobset struct {
	jobs []*Job
}

type Querier struct {
}

func NewJobset(s int) *Jobset {
	return &Jobset{jobs: make([]*Job, 0, s)}
}

func (js *Jobset) Add(job *Job) {
	js.jobs = append(js.jobs, job)
}

func (js *Jobset) Jobs() []*Job {
	return js.jobs
}

func (js *Jobset) Count() int {
	return len(js.jobs)
}

func NewQuery() *Querier {
	return &Querier{}
}

func (a *Querier) Optimize(jobs []*Job) {

}

type Result struct {
	Bitmap *idx.Bitmap
	ID     uint32
}

func (a *Querier) Exec(js *Jobset, debug bool) (roaring.Arr, error) {
	if js.Count() == 0 {
		return nil, nil
	}

	jobs := js.Jobs()

	a.Optimize(jobs)

	base := idx.NewNullBitmap()
	var ok bool
	if debug {
		println("new")
	}
	for _, j := range jobs {
		if debug {
			println(j.Name)
		}
		if j.Hash != nil {
			aid, f := j.Hash()
			if !f {
				return nil, nil
			}

			return []uint32{aid}, nil
		}

		if j.Bitmap != nil {
			base, ok = base.And(j.Bitmap())
			if !ok {
				return nil, nil
			}
		}
	}

	return base.B.BorrowArray(), nil
}
