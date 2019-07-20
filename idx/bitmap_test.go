package idx

import "testing"

func TestBitmapMerge(t *testing.T) {
	var f bool
	a := NewNullBitmap()

	b := NewBitmap()
	b.Add(1)
	b.Add(2)

	a, f = a.And(b)
	if !f {
		t.Fatal("expected true")
	}

	if a.B.String() != "{1,2}" {
		t.Fatalf("expected {1,2}, got: %s", a.B.String())
	}

	c := NewBitmap()
	c.Add(1)
	c.Add(3)

	a, f = a.And(c)

	if a.B.String() != "{1}" {
		t.Fatalf("expected {1}, got: %s", a.B.String())
	}
}
