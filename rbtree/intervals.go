/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rbtree

import "sync"

var intervalsPool = sync.Pool{
	New: func() interface{} {
		return make(Intervals, 0, 10)
	},
}

// Intervals represents a list of Intervals.
type Intervals []*Interval

// Dispose will free any consumed resources and allow this list to be
// re-allocated.
func (ivs *Intervals) Dispose() {
	for i := 0; i < len(*ivs); i++ {
		(*ivs)[i] = nil
	}

	*ivs = (*ivs)[:0]
	intervalsPool.Put(*ivs)
}

type Interval struct {
	From uint32
	To   uint32
	Id   uint32
}

func (i *Interval) Low() uint32 {
	return i.From
}

func (i *Interval) High() uint32 {
	return i.To
}

func (i *Interval) Overlaps(j *Interval) bool {
	return i.High() >= j.Low() &&
		i.Low() <= j.High()
}

func (i *Interval) ID() uint32 {
	return i.Id
}
