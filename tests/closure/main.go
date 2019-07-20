package main

func applyIndex(old uint32, new uint32, add func(), delete func(), yes func(), no func()) {
	add()
}

type st struct {
	old    uint32
	new    uint32
	add    func()
	delete func()
	yes    func()
	no     func()
}

func applyIndex2(s st) {
	s.add()
}

func main() {
	var a uint32 = 1
	var b uint32 = 2
	var c = 3
	applyIndex(
		a,
		b,
		func() {
			println(c)
		},
		func() {
		},
		nil,
		nil,
	)

	applyIndex2(
		st{
			old: a,
			new: b,
			add: func() {
				println(c)
			},
		},
	)
}
