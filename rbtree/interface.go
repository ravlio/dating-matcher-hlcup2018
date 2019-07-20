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

/*
Package augmentedtree is designed to be useful when checking
for intersection of ranges in n-dimensions.  For instance, if you imagine
an xy plane then the augmented tree is for telling you if
plane defined by the points (0, 0) and (10, 10).  The augmented
tree can tell you if that plane overlaps with a plane defined by
(-5, -5) and (5, 5) (true in this case).  You can also check
intersections against a point by constructing a range of encompassed
solely if a single point.

The current tree is a simple top-down red-black binary search tree.

TODO: AddUint32 a bottom-up implementation to assist with duplicate
range handling.
*/
package rbtree

// Tree defines the object that is returned from the
// tree constructor.  We use a Tree interface here because
// the returned tree could be a single dimension or many
// dimensions.
type Tree interface {
	// AddUint32 will add the provided intervals to the tree.
	Add(intervals ...*Interval)
	// Len returns the number of intervals in the tree.
	Len() uint64
	// DeleteUint32 will remove the provided intervals from the tree.
	Delete(intervals ...*Interval)
	// Query will return a list of intervals that intersect the provided
	// interval.  The provided interval's ID method is ignored so the
	// provided ID is irrelevant.
	Query(interval *Interval) Intervals
	// Traverse will readNode tree and give alls intervals
	// found in an undefined order
	Traverse(func(*Interval))
}
