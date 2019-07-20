package slist

import "sort"

func Insert(r *[]uint32, x uint32) {
	if len(*r) == 0 {
		*r = append(*r, x)
		return
	}
	i := sort.Search(len(*r), func(i int) bool { return (*r)[i] >= x })
	if i == len(*r) {
		*r = append(*r, x)
		return
	}
	if (*r)[i] != x {
		*r = append(*r, 0)
		copy((*r)[i+1:], (*r)[i:])
		(*r)[i] = x
	}
}

func And(a, b []uint32) (r []uint32) {
	if a == nil || b == nil {
		return nil
	}

	var s, l []uint32

	if len(a) >= len(b) {
		s, l = b, a

	} else {
		s, l = a, b
	}

	var i, j int
	si := len(s) - 1
	li := len(l) - 1
	// TODO очень хорошо подумть над оптимизацией
	for i <= si && j <= li {
		if s[i] == l[j] {
			r = append(r, s[i])
			i++
			j++
		} else if s[i] < l[j] {
			i++
		} else {
			j++
		}
	}

	return r
}
