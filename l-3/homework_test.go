package main

import (
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

/*
// Предположим, что все будут производить копирование
// буффера только с использованием метода Clone()
type COWBuffer struct { ... }

func NewCOWBuffer(data \[\]byte)                         // создать буффер с определенными данными
func (b \*COWBuffer) Clone() COWBuffer                  // создать новую копию буфера
func (b \*COWBuffer) Close()                            // перестать использовать копию буффера
func (b \*COWBuffer) Update(index int, value byte) bool // изменить определенный байт в буффере
func (b \*COWBuffer) String() string                    // сконвертировать буффер в строку
*/

func NewCOWBuffer(data []byte) COWBuffer {
	if data == nil {
		data = make([]byte, 0, 2<<6)
	}
	refs := new(int)
	*refs += 1

	return COWBuffer{data: data, refs: refs, n: len(data)}
}

var ErrRefsNil = errors.New("refs == nil")
var ErrBufNil = errors.New("buf == nil")

type COWBuffer struct {
	data  []byte
	refs  *int
	n     int
	close bool
}

func (b *COWBuffer) Clone() COWBuffer {
	if b.refs == nil {
		panic(ErrRefsNil)
	}

	*b.refs += 1
	return COWBuffer{
		refs: b.refs,
		data: b.data,
		n:    b.n,
	}
}

func (b *COWBuffer) Close() {
	if b.close {
		return
	}

	if !b.close {
		b.close = true
	}

	if b.refs == nil {
		panic(ErrRefsNil)
	}
	b.n = 0
	b.data = nil
	*b.refs -= 1
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b.refs == nil {
		panic(ErrRefsNil)
	}
	if b.data == nil {
		panic(ErrBufNil)
	}
	if index < 0 || index > len(b.data)-1 {
		return false
	}

	if *b.refs > 1 {
		*b.refs--
		b.refs = new(int)
		*b.refs += 1
		buf := make([]byte, len(b.data))
		copy(buf, b.data)
		b.data = buf
	}
	b.data[index] = value

	return true
}

func (b *COWBuffer) String() string {
	return unsafe.String(&b.data[0], b.n)
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
