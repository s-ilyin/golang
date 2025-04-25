package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Идея упорядоченного словаря заключается в том, что он будет реализован на основе бинарного дерева поиска (BST). Дерево будет строиться только по ключам элементов, значения элементов при построении дерева не учитываются. Элементы с одинаковыми ключами в упорядоченном словаре хранить нельзя.


Поподробнее с бинарными деревьями поиска можно познакомиться [здесь.](https://habr.com/ru/articles/65617/)


API для упорядоченного словаря будет выглядеть следующим образом:

type OrderedMap struct { ... }

func NewOrderedMap() OrderedMap                      // создать упорядоченный словарь
func (m \*OrderedMap) Insert(key, value int)          // добавить элемент в словарь
func (m \*OrderedMap) Erase(key int)                  // удалить элемент из словари
func (m \*OrderedMap) Contains(key int) bool          // проверить существование элемента в словаре
func (m \*OrderedMap) Size() int                      // получить количество элементов в словаре
func (m \*OrderedMap) ForEach(action func(int, int))  // применить функцию к каждому элементу словаря от меньшего к большему
*/

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

type OrderedMap struct {
	root *node
	size int
}

type node struct {
	k     int
	v     int
	left  *node
	right *node
}

// Insert добавить элемент в словарь
func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &node{k: key, v: value, left: nil, right: nil}
		m.size = 1
		return
	}
	n := m.root.insert(key, value)
	m.size += n
}

func (n *node) contains(key int) bool {
	if n == nil {
		return false
	}
	if n.k == key {
		return true
	}

	if key < n.k {
		if n.left == nil {
			return false
		}

		return n.left.contains(key)
	}

	if n.right == nil {
		return false
	}
	return n.right.contains(key)
}

func (n *node) insert(key, value int) int {
	if n.k == key {
		n.v = value
		return 0
	}
	if key < n.k {
		if n.left == nil {
			n.left = &node{k: key, v: value, left: nil, right: nil}
			return 1
		}

		n.left.insert(key, value)
		return 1
	}
	if n.right == nil {
		n.right = &node{k: key, v: value, left: nil, right: nil}
		return 1
	}
	return n.right.insert(key, value)
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (n *node) inorder(action func(int, int)) {
	if n == nil {
		return
	}
	n.left.inorder(action)
	action(n.k, n.v)
	n.right.inorder(action)
}

func (m *OrderedMap) ForEach(action func(k, v int)) {
	m.root.inorder(action)
}

func (m *OrderedMap) Contains(key int) bool {
	if m.root == nil {
		return false
	}

	return m.root.contains(key)
}

func (n *node) erase(key int) (*node, bool) {
	if key < n.k {
		if n.left == nil {
			return nil, false
		}
		left, ok := n.left.erase(key)
		n.left = left

		return n, ok
	}
	if key > n.k {
		if n.right == nil {
			return nil, false
		}

		right, ok := n.right.erase(key)
		n.right = right

		return n, ok
	}

	if n.k == key {
		if n.left == nil && n.right == nil {
			n = nil
			return n, true
		}

		if n.left == nil {
			n = n.right
			return n, true
		}
		if n.right == nil {
			n = n.left
			return n, true
		}

		small := n.right
		for small != nil && small.left != nil {
			small = small.left
		}
		n.k = small.k
		n.v = small.v
		right, ok := n.right.erase(n.k)
		n.right = right

		return n, ok
	}

	return nil, false
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	if _, ok := m.root.erase(key); ok {
		m.size--
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
