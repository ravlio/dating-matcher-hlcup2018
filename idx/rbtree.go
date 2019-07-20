package idx

import (
	"github.com/ravlio/highloadcup2018/rbtree"
	"sync"
)

type RBTree struct {
	Mx *sync.RWMutex
	T  rbtree.Tree
}

func NewRBTree() *RBTree {
	return &RBTree{T: rbtree.New()}
}
