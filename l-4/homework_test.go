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
}

type node struct {
	k     int
	v     int
	left  *node
	right *node
}

func (m *OrderedMap) Insert(key, value int) {
	m.root = insert(m.root, key, value)
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	erase(m.root, key)
}

func (m *OrderedMap) Contains(key int) bool {
	if m.root == nil {
		return false
	}

	return contains(m.root, key)
}

func (m *OrderedMap) Size() int {
	size := 0
	inorder(m.root, func(_, _ int) {
		size += 1
	})

	return size
}

func (m *OrderedMap) ForEach(action func(k, v int)) {
	inorder(m.root, action)
}

func contains(n *node, k int) bool {
	if n == nil {
		return false
	}
	if n.k == k {
		return true
	}

	if k < n.k {
		return contains(n.left, k)
	} else {
		return contains(n.right, k)
	}
}

func insert(n *node, k, v int) *node {
	if n == nil {
		return &node{k: k, v: v, left: nil, right: nil}
	}
	if n.k == k {
		n.v = v
	}

	if k < n.k {
		n.left = insert(n.left, k, v)
	}
	if k > n.k {
		n.right = insert(n.right, k, v)
	}

	return n
}

func inorder(n *node, action func(int, int)) {
	if n == nil {
		return
	}
	inorder(n.left, action)
	action(n.k, n.v)
	inorder(n.right, action)
}

func erase(n *node, k int) *node {
	if n == nil {
		return nil
	}
	if k < n.k {
		n.left = erase(n.left, k)
		return n
	}
	if k > n.k {
		n.right = erase(n.right, k)
		return n
	}

	if n.left == nil {
		return n.right
	}
	if n.right == nil {
		return n.left
	}

	min := n.right
	for min != nil && min.left != nil {
		min = min.left
	}
	n.k = min.k
	n.v = min.v
	n.right = erase(n.right, n.k)

	return n
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
