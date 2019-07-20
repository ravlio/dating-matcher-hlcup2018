// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use th file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package int32t

import (
	"errors"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/idx"
	"math"
	"math/rand"
	"reflect"
)

var (
	DefaultMaxLevel            = 12
	DefaultProbability float32 = 0.25
)

type Skiplist struct {
	// Determining MaxLevel
	// Reference: http://drum.lib.umd.edu/bitstream/1903/544/2/CS-TR-2286.1.pdf - section 2
	//
	// > To get away from magic constants, we say that a fraction p of the nodes with level i pointers
	// > also have level i+1 pointers.
	//
	// ip = inverse of p or 1/p or int(math.Ceil(1/p))
	ip int

	// > Since we can safely cap levels at L(n), we should choose MaxLevel = L(N) (where N is an upper
	// > bound on the number of elements in a skip list). If p = 1/2, using MaxLevel = 32 is appropriate
	// > for data structures containing up to 2^32 elements.
	//
	// Magic formula is L = log base 1/p of N or (1/p)^L = N
	//
	// Given p = 1/4 and L = 12, then (1/(1/4))^12 = 4^12 = 2^24 = 16777216 elements in the skiplist
	maxLevel int

	// The number of levels th list has currently. Likely increase until MaxLevel.
	// level is 0-based, so the bottom level is 0, max level is maxLevel-1
	level int

	// headNode is the first node in the skiplist. The next pointers in headNode always points forward
	// to the next node at the appropriate height. Initially all the next pointers will point to tailNode.
	// All of the prev pointers will remain nil.
	headNode *node

	// Using Search Fingers
	// Reference: http://drum.lib.umd.edu/bitstream/1903/544/2/CS-TR-2286.1.pdf - section 3.1
	// We keep two sets of fingers as search and insert localities are likely different, especially if
	// the insert keys are close to each other
	insertFingers []*node

	// fingers for selecting nodes
	selectFingers []*node

	// Total number of nodes inserted
	count int

	// Comparison function for the node keys.
	// For ascending order - if k1 < k2 return true; else return false
	// For descending order - if k1 > k2 return true; else return false
	compare Comparator
}

func New(compare Comparator) *Skiplist {
	l := DefaultMaxLevel
	ip := int(math.Ceil(1 / float64(DefaultProbability)))

	return &Skiplist{
		ip:            ip,
		maxLevel:      l,
		insertFingers: make([]*node, l),
		selectFingers: make([]*node, l),
		level:         1,
		count:         0,
		compare:       compare,
		headNode:      newNode(l),
	}
}

func (th *Skiplist) SetCompare(compare Comparator) (err error) {
	if compare == nil {
		return errors.New("skiplist/SetCompare: trying to set comparator to nil")
	}
	th.compare = compare
	return nil
}

func (th *Skiplist) SetMaxLevel(l int) (err error) {
	if l < 1 {
		return errors.New("skiplist/SetCompare: max level must be greater than zero (0)")
	}
	th.maxLevel = l
	return nil
}

func (th *Skiplist) SetProbability(p float32) (err error) {
	if p > 1 {
		p = 1
	}
	th.ip = int(math.Ceil(1 / float64(p)))
	return nil
}

func (th *Skiplist) Close() (err error) {
	return nil
}

func (th *Skiplist) Count() int {
	return th.count
}

func (th *Skiplist) Level() int {
	return th.level
}

// Choose the new node's level, branching with p (1/ip) probability, with no regards to N (size of list)
func (th *Skiplist) newNodeLevel() int {
	h := 1

	for h < th.maxLevel && rand.Intn(th.ip) == 0 {
		h++
	}

	return h
}

func (th *Skiplist) updateSearchFingers(key int32, fingers []*node) (err error) {
	startLevel := th.level - 1
	startNode := th.headNode

	if fingers[0] != nil && fingers[0].key != 0 {
		if less, err := th.compare(fingers[0].key, key); err != nil {
			return err
		} else if less {
			// Move forward, find the highest level s.t. the next node's key < key
			for l := 1; l < th.level; l++ {
				if fingers[l].next[l] != nil && fingers[l].key == 0 {
					// If the next node is not nil and fingers[l].key >= key
					if less, err := th.compare(fingers[l].key, key); err != nil {
						return err
					} else if less == false {
						startLevel = l - 1
						startNode = fingers[l]
						break
					}
				}
			}
		} else {
			//log.Println("inside if else, th.level =", th.level-1)
			// Move backward, find the lowest level s.t. the node's timestamp < t
			for l := 1; l < th.level; l++ {
				//log.Println("inside for loop, level =", l)
				// fingers[l].key < key
				if fingers[l].key != 0 {
					if less, err := th.compare(fingers[l].key, key); err != nil {
						return err
					} else if less {
						startLevel = l
						startNode = fingers[l]
						break
					}
				}
			}
		}
	}

	// For each of the skiplist levels, going from the current height to 1, walk the list until
	// we find a node that has a timestamp that's >= the timestamp t, or the end of the list
	// l = level, p = ptr to node during traversal
	for l, p := startLevel, startNode; l >= 0; l-- {
		n := p.next[l]

		for {
			if n == nil {
				// last node on the list
				// go to the next level down, and continue traversing
				//log.Println("n == nil")
				break
			}

			//log.Println("n != nil")
			// If n.key >= key
			if less, err := th.compare(n.key, key); err != nil {
				return err
			} else if less == false {
				// Found the first record that either has the same timestamp or greater at th level
				// go to the next level down, and continue traversing
				//log.Println("nt >= t, nt = ", nt.(int64))
				break
			}
			//log.Println("after compare")

			// Move the pointers forward, p = n, n = n.next
			p, n = n, n.next[l]
		}

		fingers[l] = p
	}

	return nil
}

func (th *Skiplist) Insert(key int32, value uint32) (*node, error) {
	if key == 0 {
		return nil, errors.New("skiplist/Insert: key is nil")
	}

	if th.compare == nil {
		return nil, errors.New("skiplist/Insert: comparator is not set (== nil)")
	}

	// Create new node
	l := th.newNodeLevel()
	n := newNode(l)
	n.SetKey(key)
	n.SetValue(value)

	//log.Println("th.finger[0] =", th.insertFingers[0])
	// Find the position where we should insert the node by updating the search insertFingers using the key
	// Search insertFingers will be updated with the rightmost element of each level that is left of the element
	// that's greater than or equal to key.
	// In other words, we are inserting the new node to the right of the search insertFingers.
	if err := th.updateSearchFingers(key, th.insertFingers); err != nil {
		return nil, errors.New("skiplist/insert: cannot find insert position, " + err.Error())
	}

	//log.Println("search insertFingers =", th.insertFingers)
	// Raise the level of the skiplist if the new level is higher than the existing list level
	// So for levels higher than the current list level, the previous node is headNode for that level
	if th.level < l {
		for i := th.level; i < l; i++ {
			//log.Println("before ---- ", th.insertFingers)
			th.insertFingers[i] = th.headNode
			//log.Println("after  ---- ", th.insertFingers)
		}
		th.level = l
		//log.Println("new th.level =", l)
		//log.Println("th.insertFingers =", th.insertFingers)
	}

	// Finally insert the node into the skiplist
	for i := 0; i < l; i++ {
		// new node points forward to the previous node's next node
		// previous node's next node points to the new node
		n.next[i], th.insertFingers[i].next[i] = th.insertFingers[i].next[i], n
	}

	// Adding to the count
	th.count++

	return n, nil
}

// Select a list of nodes that match the key. The results are stored in the array pointed to by results
func (th *Skiplist) Select(key int32) (iter *Iterator, err error) {
	return th.SelectRange(key, key)
}

func (th *Skiplist) SelectFromBitmap(key1 int32) (bm *idx.Bitmap, found bool, err error) {
	if err = th.updateSearchFingers(key1, th.selectFingers); err != nil {
		return nil, false, errors.New("skiplist/SelectRange: error selecting nodes, " + err.Error())
	}

	f := false
	bm = idx.NewBitmap()
	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		f = true
		bm.AddUnsafe(p.value)

	}

	if f {
		return bm, f, nil

	} else {
		return nil, false, nil
	}
}

func (th *Skiplist) SelectRange(key1, key2 int32) (iter *Iterator, err error) {
	if key1 == 0 || key2 == 0 {
		return nil, errors.New("skiplist/SelectRange: key1 or key2 is nil")
	}

	if th.compare == nil {
		return nil, errors.New("skiplist/SelectRange: comparator is not set (== nil)")
	}

	// Walk the levels and nodes until we find the node at the lowest level (0) that the comparator returns false
	// E.g., if comparator is BuiltinLessThan, then we find the node at the lowest level s.t. node.key < key
	// Then we walk from there to find all the nodes that have node.key == key
	// We keep track of the last touched nodes at each level as selectFingers, and then we re-use the selectFingers
	// so that we can get O(log k) where k is the distance between last searched key and current search key
	// -- ok, so all th is done by updateSearchFingers

	if err = th.updateSearchFingers(key1, th.selectFingers); err != nil {
		return nil, errors.New("skiplist/SelectRange: error selecting nodes, " + err.Error())
	}

	iter = newIterator()
	var res bool
	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		pk := p.GetKey()
		if res, err = th.compare(pk, key2); err != nil {
			// If there's error in comparing the keys, then return err
			return nil, errors.New("skiplist/SelectRange: error comparing keys; " + err.Error())
		} else if res || reflect.DeepEqual(pk, key2) {
			iter.buf = append(iter.buf, p)
			iter.count++
		} else {
			// Otherwise if the p.key is "after" key, after could mean greater or less, depending
			// on the comparator, then we know we are done
			break
		}
	}

	return iter, nil
}

func (th *Skiplist) SelectFrom(key1 int32) (iter *Iterator, err error) {
	if key1 == 0 {
		return nil, errors.New("skiplist/SelectRange: key1 or key2 is nil")
	}

	if th.compare == nil {
		return nil, errors.New("skiplist/SelectRange: comparator is not set (== nil)")
	}

	// Walk the levels and nodes until we find the node at the lowest level (0) that the comparator returns false
	// E.g., if comparator is BuiltinLessThan, then we find the node at the lowest level s.t. node.key < key
	// Then we walk from there to find all the nodes that have node.key == key
	// We keep track of the last touched nodes at each level as selectFingers, and then we re-use the selectFingers
	// so that we can get O(log k) where k is the distance between last searched key and current search key
	// -- ok, so all th is done by updateSearchFingers

	if err = th.updateSearchFingers(key1, th.selectFingers); err != nil {
		return nil, errors.New("skiplist/SelectRange: error selecting nodes, " + err.Error())
	}

	iter = newIterator()
	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		iter.buf = append(iter.buf, p)
		iter.count++
	}

	return iter, nil
}

func (th *Skiplist) SearchFrom(key1 int32, limit int) (bm *roaring.Bitmap, err error) {
	if key1 == 0 {
		return nil, errors.New("skiplist/SelectRange: key1 or key2 is nil")
	}

	if th.compare == nil {
		return nil, errors.New("skiplist/SelectRange: comparator is not set (== nil)")
	}

	// Walk the levels and nodes until we find the node at the lowest level (0) that the comparator returns false
	// E.g., if comparator is BuiltinLessThan, then we find the node at the lowest level s.t. node.key < key
	// Then we walk from there to find all the nodes that have node.key == key
	// We keep track of the last touched nodes at each level as selectFingers, and then we re-use the selectFingers
	// so that we can get O(log k) where k is the distance between last searched key and current search key
	// -- ok, so all th is done by updateSearchFingers

	if err = th.updateSearchFingers(key1, th.selectFingers); err != nil {
		return nil, errors.New("skiplist/SelectRange: error selecting nodes, " + err.Error())
	}

	var i = 0
	bm = roaring.New()
	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		if limit > 0 && i > limit {
			break
		}
		bm.Add(p.value)

	}

	return bm, nil
}

func (th *Skiplist) Delete(key int32) (iter *Iterator, err error) {
	return th.DeleteRange(key, key)
}

func (th *Skiplist) DeleteRange(key1, key2 int32) (iter *Iterator, err error) {
	if key1 == 0 || key2 == 0 {
		return nil, errors.New("skiplist/DeleteRange: key1 or key2 is nil")
	}

	if reflect.TypeOf(key1) != reflect.TypeOf(key2) {
		return nil, fmt.Errorf("skiplist/DeleteRange: k1.(%s) and k2.(%s) have different types",
			reflect.TypeOf(key1).Name(), reflect.TypeOf(key2).Name())
	}

	if th.compare == nil {
		return nil, errors.New("skiplist/DeleteRange: comparator is not set (== nil)")
	}

	// Walk the levels and nodes until we find the node at the lowest level (0) that the comparator returns false
	// E.g., if comparator is BuiltinLessThan, then we find the node at the lowest level s.t. node.key < key
	// Then we walk from there to find all the nodes that have node.key == key
	// We keep track of the last touched nodes at each level as selectFingers, and then we re-use the selectFingers
	// so that we can get O(log k) where k is the distance between last searched key and current search key
	// -- ok, so all th is done by updateSearchFingers

	if err = th.updateSearchFingers(key1, th.selectFingers); err != nil {
		return nil, errors.New("skiplist/DeleteRange: error finding node; " + err.Error())
	}

	iter = newIterator()
	var res bool
	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		pk := p.GetKey()
		if res, err = th.compare(pk, key2); err != nil {
			// If there's error in comparing the keys, then return err
			return nil, errors.New("skiplist/DeleteRange: error comparing keys; " + err.Error())
		} else if res || reflect.DeepEqual(pk, key2) {
			iter.buf = append(iter.buf, p)
			iter.count++

			for i := 0; i < th.level; i++ {
				if th.selectFingers[i].next[i] != p {
					break
				}
				th.selectFingers[i].next[i] = p.next[i]
			}

			th.count--

			for th.level > 1 && th.headNode.next[th.level-1] == nil {
				th.level--
			}
		} else {
			// Otherwise if the p.key is "after" key, after could mean greater or less, depending
			// on the comparator, then we know we are done
			break
		}
	}

	return iter, nil
}

func (th *Skiplist) RealCount(i int) (c int) {
	for p := th.headNode.next[i]; p != nil; {
		if p != nil {
			//log.Println("node =", p.record)
			c++
			p = p.next[i]
		}
	}

	return
}

func (th *Skiplist) PrintStats() {
	fmt.Println("Real count   :", th.RealCount(0))
	fmt.Println("Total levels :", th.Level())

	for i := 0; i < th.level; i++ {
		fmt.Println("Level", i, "count:", th.RealCount(i))
	}

}

func (th *Skiplist) DeleteValue(key int32, v uint32) bool {
	if err := th.updateSearchFingers(key, th.selectFingers); err != nil {
		return false
	}

	for p := th.selectFingers[0].next[0]; p != nil; p = p.next[0] {
		if p.GetKey() == key && p.GetValue() == v {
			for i := 0; i < th.level; i++ {
				if th.selectFingers[i].next[i] != p {
					break
				}
				th.selectFingers[i].next[i] = p.next[i]
			}

			th.count--

			for th.level > 1 && th.headNode.next[th.level-1] == nil {
				th.level--
			}
			return true
		}

	}

	return false
}
