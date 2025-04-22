package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

/*
CircularQueue API

func (q \*CircularQueue) Push(value int) bool // добавить значение в конец очереди (false, если очередь заполнена)
func (q \*CircularQueue) Pop() bool           // удалить значение из начала очереди (false, если очередь пустая)
func (q \*CircularQueue) Front() int          // получить значение из начала очереди (-1, если очередь пустая)
func (q \*CircularQueue) Back() int           // получить значение из конца очереди (-1, если очередь пустая)
func (q \*CircularQueue) Empty() bool         // проверить пустая ли очередь
func (q \*CircularQueue) Full() bool          // проверить заполнена ли очередь
*/

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{values: make([]int, size), head: -1, tail: -1}
}

type CircularQueue struct {
	values     []int
	head, tail int // head - read idx, head - write idx
}

func (q *CircularQueue) Full() bool {
	return q.full()
}

func (q *CircularQueue) Empty() bool {
	return q.emptyHead() && q.emptyTail()
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}

	return q.values[q.head]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.tail]
}

func (q *CircularQueue) Push(value int) bool {
	if q.full() {
		return false
	}
	if q.head == -1 {
		q.head = 0
	}
	q.tail = (q.tail + 1) % len(q.values)
	q.values[q.tail] = value

	return true
}

func (q *CircularQueue) full() bool {
	return q.head == (q.tail+1)%len(q.values)
}

func (q *CircularQueue) equal() bool {
	return q.head == q.tail
}

func (q *CircularQueue) emptyHead() bool {
	return q.head == -1
}

func (q *CircularQueue) emptyTail() bool {
	return q.tail == -1
}

func (q *CircularQueue) Pop() bool {
	if q.equal() && q.emptyHead() {
		return false
	}
	if q.equal() && !q.emptyHead() {
		_ = q.values[q.head]
		q.head = -1
		q.tail = -1
		return true
	}

	_ = q.values[q.head]
	q.head = (q.head + 1) % len(q.values)
	return true
}

func Test_Push(t *testing.T) {
	tests := []struct {
		name string
		val  int
		push bool
		full bool
	}{
		{
			name: "idx_0",
			val:  1,
			push: true,
			full: false,
		},
		{
			name: "idx_1",
			val:  2,
			push: true,
			full: false,
		},
		{
			name: "idx_2",
			val:  3,
			push: true,
			full: false,
		},
		{
			name: "idx_3",
			val:  3,
			push: true,
			full: true,
		},
		{
			name: "idx_0",
			val:  4,
			push: false,
			full: true,
		},
	}
	cq := NewCircularQueue(4)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cq.Push(tt.val); tt.push != got {
				t.Errorf("push: expect %t != got %t\n", tt.push, got)
			}
			if got := cq.full(); tt.full != got {
				t.Errorf("full: expect %t != got %t\n", tt.full, got)
			}
		})
	}
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty()) // ok
	assert.False(t, queue.Full()) // ok

	assert.Equal(t, -1, queue.Front()) // ok
	assert.Equal(t, -1, queue.Back())  // ok
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1)) // [1] back 0
	assert.False(t, queue.Empty())
	assert.True(t, queue.Push(2)) // [1, 2] back 1
	assert.False(t, queue.Empty())
	assert.True(t, queue.Push(3))  // [1, 2, 3] back 2
	assert.False(t, queue.Push(4)) // [1, 2, 3] back 2

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty()) // ok
	assert.True(t, queue.Full())   // ok

	assert.Equal(t, 1, queue.Front()) // ok front 0
	assert.Equal(t, 3, queue.Back())  // ok back 2

	assert.True(t, queue.Pop())    // del 1 [2, 3] front 1
	assert.False(t, queue.Empty()) // ok front 1 back 2
	assert.False(t, queue.Full())  // ok // front 1 back 2
	assert.True(t, queue.Push(4))  // front 1 back 0

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front()) // front 1 back 0
	assert.Equal(t, 4, queue.Back())  // front 1 back 0

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.True(t, queue.Pop())  // front 1 back 0
	assert.True(t, queue.Pop())  // front 2 back 0
	assert.True(t, queue.Pop())  // front -1 back -1
	assert.False(t, queue.Pop()) //  front -1 back -1

	assert.True(t, queue.Empty()) // front -1 back -1
	assert.False(t, queue.Full()) // front -1 back -1
}

// init
// f, r = -1, -1
// [ 0 , 0 , 0 ]

// push 1
// curr f = 0
// curr r = 0
//  f,r
//  | |
// [ 1 , 0 , 0 ]
// next f = 0
// next r = 1

// push 2
// curr f = 0
// curr r = 1
//   f   r
//   |   |
// [ 1 , 2 , 0 ]
// next f = 0
// next r = 2

// push 3
// curr f = 0
// curr r = 2
//   f       r
//   |       |
// [ 1 , 2 , 3 ]
// next f = 0
// next r = 0

// pop
// curr f = 0
// curr r = 0
//   r   f
//   |   |
// [ 0 , 2 , 3 ]
// next f = 1
// next r = 0

// pop
// curr f = 1
// curr r = 0
//   r   f
//   |   |
// [ 0 , 0 , 3 ]
// next f = 2
// next r = 0

// pop
// curr f = 2
// curr r = 0
//   r       f
//   |       |
// [ 0 , 0 , 0 ]
// next f = 0
// next r = 0

// push 4
//  f r
//  | |
// [ 4 , 0 , 0 ]
// f = 3
// r = 1
