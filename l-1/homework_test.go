package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian(number uint32) uint32 {

	// сдвигам, получаем ячейки байт
	var b1, b2, b3, b4 = number >> 0, number >> 8, number >> 16, number >> 24

	// сбрасываем и складываем
	return ((0xFF & b1) << 24) | ((0xFF & b2) << 16) | ((0xFF & b3) << 8) | ((0xFF & b4) << 0) // need to implement
}

func TestСonversion(t *testing.T) {
	// binary.BigEndian.Uint32()
	// binary.LittleEndian.Uint32()

	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #6": {
			number: 0x000000FF,
			result: 0xFF000000,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
