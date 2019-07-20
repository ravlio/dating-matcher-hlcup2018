package main

func main() {
	id := 1
	s := make([]struct{}, 3)
	s = append(s[:id-1], s[id+1:]...)

	for k := range s {
		println(k + 1)
	}

}
