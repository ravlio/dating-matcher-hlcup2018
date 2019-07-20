package indexer

type VarType int8

const (
	VarUint32 VarType = iota
	VarInt32
	VarInt64
	VarUint32Slice
	VarBool
)

// Индексатор и процессор условий добавления/удаления индексов
// По хорошему, чтобы сэкономить на колбэках,
// условия лучше заинлайнить, но руками очень сложно не ошибиться.
// Возможно тут помог бы генератор кода
// Плюс такой подход позволяет паралелить индексацию

type Job struct {
	VarType           VarType
	CurUint32         uint32
	NewUint32         uint32
	CurInt64          int64
	NewInt64          int64
	CurInt32          int32
	NewInt32          int32
	CurUint32Slice    []uint32
	NewUint32Slice    []uint32
	NewStr            string
	CurCond           bool
	NewCond           bool
	Add               func() error
	Delete            func() error
	AddUint32         func(id uint32) error
	DeleteUint32      func(id uint32) error
	AddInt64          func(id int64) error
	DeleteInt64       func(id int64) error
	AddInt32          func(id int32) error
	DeleteInt32       func(id int32) error
	AddUint32Slice    func(id []uint32) error
	DeleteUint32Slice func(id []uint32) error
	AddYes            func() error
	DeleteYes         func() error
	AddNo             func() error
	DeleteNo          func() error
}

type Jobset struct {
	jobs []*Job
}

type Indexer struct {
}

func NewJobset() *Jobset {
	return &Jobset{jobs: make([]*Job, 0)}
}

func (js *Jobset) Add(job *Job) {
	js.jobs = append(js.jobs, job)
}

func (js *Jobset) Jobs() []*Job {
	return js.jobs
}

func (i *Job) Eq() bool {
	switch i.VarType {
	case VarUint32:
		return i.CurUint32 == i.NewUint32

	case VarInt32:
		return i.CurInt32 == i.NewInt32

	case VarInt64:
		return i.CurInt64 == i.NewInt64
	case VarBool:
		return i.CurCond == i.NewCond
	}

	return false
}

func (i *Job) IsCur() bool {
	switch i.VarType {
	case VarUint32:
		return i.CurUint32 > 0
	case VarUint32Slice:
		return i.CurUint32Slice != nil

	case VarInt32:
		return i.CurInt32 > 0

	case VarInt64:
		return i.CurInt64 > 0
	}

	return false
}

func (i *Job) IsNew() bool {
	switch i.VarType {
	case VarUint32:
		return i.NewUint32 > 0

	case VarUint32Slice:
		return i.NewUint32Slice != nil

	case VarInt32:
		return i.NewInt32 > 0

	case VarInt64:
		return i.NewInt64 > 0
	}

	return false
}

func (i *Job) addCur() error {
	switch i.VarType {
	case VarUint32:
		if i.AddUint32 != nil {
			return i.AddUint32(i.CurUint32)
		}
	case VarUint32Slice:
		if i.AddUint32 != nil {
			for _, v := range i.CurUint32Slice {
				i.AddUint32(v)
			}
		}
	case VarInt32:
		if i.AddInt32 != nil {
			return i.AddInt32(i.CurInt32)
		}

	case VarInt64:
		if i.AddInt64 != nil {
			return i.AddInt64(i.CurInt64)
		}
	}

	if i.Add != nil {
		return i.Add()
	}

	return nil
}

func (i *Job) deleteCur() error {
	switch i.VarType {
	case VarUint32:
		if i.DeleteUint32 != nil {
			return i.DeleteUint32(i.CurUint32)
		}
	case VarUint32Slice:
		if i.AddUint32 != nil {
			// TODO оптимизироваться. Если идёт изменение,
			//  то старые ид удаляются все, даже если они в последствии добавятся. Нужно делать вначале пересечение
			// старых и новых слайсов
			for _, v := range i.CurUint32Slice {
				i.DeleteUint32(v)
			}
		}

	case VarInt32:
		if i.DeleteInt32 != nil {
			return i.DeleteInt32(i.CurInt32)
		}

	case VarInt64:
		if i.DeleteInt64 != nil {
			return i.DeleteInt64(i.CurInt64)
		}
	}

	if i.Delete != nil {
		return i.Delete()
	}

	return nil
}

func (i *Job) addNew() error {
	switch i.VarType {
	case VarUint32:
		if i.AddUint32 != nil {
			return i.AddUint32(i.NewUint32)
		}
	case VarUint32Slice:
		if i.AddUint32 != nil {
			for _, v := range i.NewUint32Slice {
				i.AddUint32(v)
			}
		}

	case VarInt32:
		if i.AddInt32 != nil {
			return i.AddInt32(i.NewInt32)
		}

	case VarInt64:
		if i.AddInt64 != nil {
			return i.AddInt64(i.NewInt64)
		}
	}

	if i.Add != nil {
		return i.Add()
	}

	return nil
}

func (i *Job) deleteNew() error {
	switch i.VarType {
	case VarUint32:
		if i.DeleteUint32 != nil {
			return i.DeleteUint32(i.NewUint32)
		}
	case VarUint32Slice:
		if i.AddUint32 != nil {
			for _, v := range i.NewUint32Slice {
				i.DeleteUint32(v)
			}
		}

	case VarInt32:
		if i.DeleteInt32 != nil {
			return i.DeleteInt32(i.NewInt32)
		}

	case VarInt64:
		if i.DeleteInt64 != nil {
			return i.DeleteInt64(i.NewInt64)
		}
	}

	if i.Delete != nil {
		return i.Delete()
	}

	return nil
}

func (i *Job) addYes() error {
	if i.AddYes != nil {
		return i.AddYes()
	}

	return nil
}

func (i *Job) deleteYes() error {
	if i.DeleteYes != nil {
		return i.DeleteYes()
	}

	return nil
}

func (i *Job) addNo() error {
	if i.AddNo != nil {
		return i.AddNo()
	}

	return nil
}

func (i *Job) deleteNo() error {
	if i.DeleteNo != nil {
		return i.DeleteNo()
	}

	return nil
}

func NewIndexer() *Indexer {
	return &Indexer{}
}

func (a *Indexer) Optimize(jobs []*Job) {

}
func (a *Indexer) Insert(js *Jobset) error {
	return a.build(js.Jobs(), false)
}

func (a *Indexer) Update(js *Jobset) error {
	return a.build(js.Jobs(), true)
}

func (a *Indexer) build(j []*Job, update bool) error {
	for _, v := range j {
		if update {
			if v.VarType == VarBool && !v.Eq() {
				if v.CurCond {
					if err := v.deleteCur(); err != nil {
						return err
					}
					if err := v.deleteYes(); err != nil {
						return err
					}

					if err := v.addNo(); err != nil {
						return err
					}
				} else {
					if err := v.addNew(); err != nil {
						return err
					}
					if err := v.addYes(); err != nil {
						return err
					}

					if err := v.deleteNo(); err != nil {
						return err
					}
				}

				continue
			}

			if v.Eq() {
				continue
			}

			isCur := v.IsCur()
			if isCur { // есть текущее значение. Убираем его в пользу нового
				if err := v.deleteCur(); err != nil {
					return err
				}
			}

			isNew := v.IsNew()
			if isNew { // есть новое значение. Добавляем его
				if err := v.addNew(); err != nil {
					return err
				}
			}

			if !isCur && isNew {
				if err := v.addYes(); err != nil {
					return err
				}
				if err := v.deleteNo(); err != nil {
					return err
				}
			}

			if isCur && !isNew {
				if err := v.deleteYes(); err != nil {
					return err
				}
				if err := v.addNo(); err != nil {
					return err
				}
			}

			continue
		}
		if v.VarType == VarBool {
			if v.NewCond {
				if err := v.addNew(); err != nil {
					return err
				}
				if err := v.addYes(); err != nil {
					return err
				}
			} else {
				if err := v.addNo(); err != nil {
					return err
				}
			}

			continue
		}

		if v.IsNew() {
			if err := v.addNew(); err != nil {
				return err
			}
			if err := v.addYes(); err != nil {
				return err
			}
		} else {
			if err := v.addNo(); err != nil {
				return err
			}
		}

		continue
	}

	return nil
}
