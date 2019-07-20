package indexer

import "testing"

func TestInsert0(t *testing.T) {
	idx := NewIndexer()
	var no bool

	js := NewJobset()
	js.Add(&Job{
		VarType:  VarInt32,
		NewInt32: 0,
		AddInt32: func(id int32) error {
			t.Errorf("unwanted AddInt32 call")

			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted DeleteInt32 call")
			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			no = true
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Insert(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !no {
		t.Error("missing AddNo call")
	}
}

func TestInsert1(t *testing.T) {
	idx := NewIndexer()
	var yes, add bool

	js := NewJobset()
	js.Add(&Job{
		VarType:  VarInt32,
		CurInt32: 4,
		NewInt32: 1,
		AddInt32: func(id int32) error {
			if id != 1 {
				t.Errorf("AddInt32 error. Got: %d, whant: 1", id)
			}

			add = true
			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted Delete call")
			return nil
		},
		AddYes: func() error {
			yes = true
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Insert(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !yes {
		t.Error("missing AddYes call")
	}

	if !add {
		t.Error("missing AddInt32 call")
	}
}

func TestUpdateFrom0To1(t *testing.T) {
	idx := NewIndexer()
	var yes, no, add bool
	js := NewJobset()
	js.Add(&Job{
		VarType:  VarInt32,
		CurInt32: 0,
		NewInt32: 1,
		AddInt32: func(id int32) error {
			if id != 1 {
				t.Errorf("AddInt32 error. Got: %d, whant: 1", id)
			}

			add = true
			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted Delete call")
			return nil
		},
		AddYes: func() error {
			yes = true
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			no = true
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !yes {
		t.Error("missing AddYes call")
	}

	if !no {
		t.Error("missing DeleteNo call")
	}

	if !add {
		t.Error("missing AddInt32 call")
	}
}

func TestUpdateFrom1To2(t *testing.T) {
	idx := NewIndexer()
	var add, del bool
	js := NewJobset()
	js.Add(&Job{
		VarType:  VarInt32,
		CurInt32: 1,
		NewInt32: 2,
		AddInt32: func(id int32) error {
			if id != 2 {
				t.Errorf("AddInt32 error. Got: %d, whant: 2", id)
			}

			add = true
			return nil
		},
		DeleteInt32: func(id int32) error {
			if id != 1 {
				t.Errorf("DeleteInt32 error. Got: %d, whant: 1", id)
			}

			del = true

			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !del {
		t.Error("missing DeleteInt32 call")
	}

	if !add {
		t.Error("missing AddInt32 call")
	}
}

func TestUpdateFrom2To0(t *testing.T) {
	idx := NewIndexer()
	var del, yes, no bool
	js := NewJobset()
	js.Add(&Job{
		VarType:  VarInt32,
		CurInt32: 2,
		NewInt32: 0,
		AddInt32: func(id int32) error {
			t.Errorf("unwanted AddInt32 call")
			return nil
		},
		DeleteInt32: func(id int32) error {
			if id != 2 {
				t.Errorf("DeleteInt32 error. Got: %d, whant: 2", id)
			}

			del = true

			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			no = true
			return nil
		},
		DeleteYes: func() error {
			yes = true
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !del {
		t.Error("missing DeleteInt32 call")
	}

	if !yes {
		t.Error("missing DeleteYes call")
	}

	if !no {
		t.Error("missing AddNo call")
	}
}

/////////
// bool
////////

func TestInsertBoolFalse(t *testing.T) {
	idx := NewIndexer()
	var no bool
	js := NewJobset()
	js.Add(&Job{
		VarType:  VarBool,
		NewCond:  false,
		NewInt32: 0,
		AddInt32: func(id int32) error {
			t.Errorf("unwanted AddInt32 call")

			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted DeleteInt32 call")
			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			no = true
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Insert(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !no {
		t.Error("missing AddNo call")
	}
}

func TestInsertBoolTrue(t *testing.T) {
	idx := NewIndexer()
	var yes, add bool

	js := NewJobset()
	js.Add(&Job{
		VarType: VarBool,
		NewCond: true,
		Add: func() error {
			add = true

			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted Delete call")
			return nil
		},
		AddYes: func() error {
			yes = true
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Insert(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !yes {
		t.Error("missing AddYes call")
	}

	if !add {
		t.Error("missing AddInt32 call")
	}
}

func TestUpdateBoolFromFalseToTrue(t *testing.T) {
	idx := NewIndexer()
	var yes, no, add bool

	js := NewJobset()
	js.Add(&Job{
		VarType: VarBool,
		CurCond: false,
		NewCond: true,
		Add: func() error {
			add = true
			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted Delete call")
			return nil
		},
		AddYes: func() error {
			yes = true
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			no = true
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !yes {
		t.Error("missing AddYes call")
	}

	if !no {
		t.Error("missing DeleteNo call")
	}

	if !add {
		t.Error("missing AddInt32 call")
	}
}

func TestUpdateBoolFromTrueToTrue(t *testing.T) {
	idx := NewIndexer()

	js := NewJobset()
	js.Add(&Job{
		VarType: VarBool,
		CurCond: true,
		NewCond: true,
		Add: func() error {
			t.Errorf("unwanted Add call")
			return nil
		},
		Delete: func() error {
			t.Errorf("unwanted Delete call")

			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			t.Errorf("unwanted AddNo call")
			return nil
		},
		DeleteYes: func() error {
			t.Errorf("unwanted DeleteYes call")
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}

}

func TestUpdateFromTrueToFalse(t *testing.T) {
	idx := NewIndexer()
	var del, yes, no bool
	js := NewJobset()
	js.Add(&Job{
		VarType: VarBool,
		CurCond: true,
		NewCond: false,
		Add: func() error {
			t.Errorf("unwanted Add call")
			return nil
		},
		Delete: func() error {
			del = true

			return nil
		},
		AddYes: func() error {
			t.Errorf("unwanted AddYes call")
			return nil
		},
		AddNo: func() error {
			no = true
			return nil
		},
		DeleteYes: func() error {
			yes = true
			return nil
		},
		DeleteNo: func() error {
			t.Errorf("unwanted DeleteNo call")
			return nil
		},
	})

	err := idx.Update(js)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !del {
		t.Error("missing DeleteInt32 call")
	}

	if !yes {
		t.Error("missing DeleteYes call")
	}

	if !no {
		t.Error("missing AddNo call")
	}
}
