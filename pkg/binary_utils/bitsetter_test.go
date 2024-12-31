package binary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitSetter_Set(t *testing.T) {
	tests := []struct {
		name     string
		initial  byte
		idx      uint8
		value    bool
		expected byte
	}{
		{
			name:     "Set bit 0 (leftmost) to 1",
			initial:  0b00000000,
			idx:      0,
			value:    true,
			expected: 0b10000000,
		},
		{
			name:     "Set bit 7th to 1",
			initial:  0b00000000,
			idx:      7,
			value:    true,
			expected: 0b00000001,
		},
		{
			name:     "Clear 3th bit ",
			initial:  0b00010000,
			idx:      3,
			value:    false,
			expected: 0b00000000,
		},
		{
			name:     "Set bit 5 (6th bit from left) to 1 when already 1",
			initial:  0b00100000,
			idx:      5,
			value:    true,
			expected: 0b00100100,
		},
		{
			name:     "Clear bit 3 (4th bit from left) to 0",
			initial:  0b00010000,
			idx:      3,
			value:    false,
			expected: 0b00000000,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := NewBitSetter(tc.initial)
			bs.Set(tc.idx, tc.value)
			assert.Equal(t, tc.expected, bs.Byte())
		})
	}
}

func TestBitSetter_SetRange(t *testing.T) {
	tests := []struct {
		name       string
		initial    byte
		startIndex uint8
		endIndex   uint8
		value      uint8
		expected   byte
	}{
		{
			name:       "Set range [0, 3] (leftmost 4 bits) to 1010",
			initial:    0b00000000,
			startIndex: 0,
			endIndex:   3,
			value:      0b1010,
			expected:   0b10100000,
		},
		{
			name:       "Set range [3, 7] (bits 4 to 7 from left) to 1110",
			initial:    0b00000000,
			startIndex: 4,
			endIndex:   7,
			value:      0b1110,
			expected:   0b00001110,
		},
		{
			name:       "Set range [0, 7] (all bits) to 11001100",
			initial:    0b00000000,
			startIndex: 0,
			endIndex:   7,
			value:      0b11001100,
			expected:   0b11001100,
		},
		{
			name:       "Clear range [0, 3] (leftmost 4 bits) to 0000",
			initial:    0b11110000,
			startIndex: 0,
			endIndex:   3,
			value:      0b0000,
			expected:   0b00000000,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := NewBitSetter(tc.initial)
			bs.SetRange(tc.startIndex, tc.endIndex, tc.value)
			assert.Equal(t, tc.expected, bs.Byte())
		})
	}
}

func TestToBinaryArray(t *testing.T) {
	tests := []struct {
		name     string
		value    uint8
		expected []bool
	}{
		{
			name:     "Convert 0 to binary array",
			value:    0,
			expected: []bool{},
		},
		{
			name:     "Convert 1 to binary array",
			value:    1,
			expected: []bool{true},
		},
		{
			name:     "Convert 2 to binary array",
			value:    2,
			expected: []bool{true, false},
		},
		{
			name:     "Convert 5 to binary array",
			value:    5,
			expected: []bool{true, false, true},
		},
		{
			name:     "Convert 8 to binary array",
			value:    8,
			expected: []bool{true, false, false, false},
		},
		{
			name:     "Convert 255 to binary array",
			value:    255,
			expected: []bool{true, true, true, true, true, true, true, true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := toBinaryArray(tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}
