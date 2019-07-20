package trie

import (
	"fmt"
	"github.com/ravlio/highloadcup2018/idx"
	"sync"
)
import "github.com/ravlio/highloadcup2018/trie/patricia"

type Trie struct {
	Mx sync.RWMutex
	T  *patricia.Trie
}

func (t *Trie) Insert(k string, uid uint32) {
	t.Mx.Lock()
	// Добавляем строку и в конец - ид пользователя, чтобы строка в дереве была не уникальной
	// Так как две одинаковые страны добавить нельзя
	// TODO переделать на быстрое решение с байтами
	t.T.Insert([]byte(fmt.Sprintf("%s$%d", k, uid)), patricia.Item(uid))
	t.Mx.Unlock()
}

func (t *Trie) Delete(k string, uid uint32) {
	t.T.Delete([]byte(fmt.Sprintf("%s$%d", k, uid)))
}

func (t *Trie) Select(k string) *idx.Bitmap {
	b := idx.NewBitmap()
	t.T.VisitSubtree(patricia.Prefix(k), func(prefix patricia.Prefix, item patricia.Item) error {
		b.AddUnsafe(uint32(item))
		return nil
	})

	return b
}

func NewTrie() *Trie {
	return &Trie{T: patricia.NewTrie()}
}
