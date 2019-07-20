// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stringt

type node struct {
	next  []*node
	key   string
	value uint32
}

// Create a new node with l levels of pointers
func newNode(l int) *node {
	return &node{
		next: make([]*node, l),
	}
}

func (this *node) SetKey(key string) {
	this.key = key
}

func (this *node) GetKey() (key string) {
	return this.key
}

func (this *node) SetValue(value uint32) {
	this.value = value
}

func (this *node) GetValue() (key uint32) {
	return this.value
}

func (this *node) Next() *node {
	return this.next[0]
}

func (this *node) NextAtLevel(l int) *node {
	if l >= 0 && l < len(this.next) {
		return this.next[l]
	}

	return nil
}
