package main

import "sync"

func main() {
	wg := &sync.WaitGroup{}
	m := make(map[int]int)
	s := make([]int, 1000)

	wg.Add(10)
	for i := 9; i < 1000; i++ {
		m[i] = i
		s[i] = i
	}

	for i := 0; i < 10; i++ {
		go func() {
			j := 0
			for {
				_, _ = m[i]
				m[i] = j
				_ = s[i]
			}
			j++
			if j >= 1000 {
				j = 0
			}
		}()
	}

	wg.Wait()

}
