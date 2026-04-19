package huffman

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "empty string returns error",
			input:     "",
			assertion: assert.Error,
		},
		{
			name:      "single character",
			input:     "a",
			assertion: assert.NoError,
		},
		{
			name:      "repeated single character",
			input:     "aaaa",
			assertion: assert.NoError,
		},
		{
			name:      "two distinct characters",
			input:     "ab",
			assertion: assert.NoError,
		},
		{
			name:      "short word",
			input:     "hello",
			assertion: assert.NoError,
		},
		{
			name:      "sentence with spaces",
			input:     "hello world",
			assertion: assert.NoError,
		},
		{
			name:      "long repeated pattern",
			input:     strings.Repeat("abcde", 100),
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := compress(tt.input)
			tt.assertion(t, err, "compress(%q)", tt.input)
		})
	}
}

func TestWriteBits(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		codes        map[rune]string
		expectedBits []byte
	}{
		{
			name:         "single character single bit",
			text:         "a",
			codes:        map[rune]string{'a': "0"},
			expectedBits: []byte{0b00000000},
		},
		{
			name:         "two characters filling one byte",
			text:         "abababab",
			codes:        map[rune]string{'a': "0", 'b': "1"},
			expectedBits: []byte{0b01010101},
		},
		{
			name:  "bits padded to fill last byte",
			text:  "ab",
			codes: map[rune]string{'a': "0", 'b': "1"},
			// "01" padded to "01000000"
			expectedBits: []byte{0b01000000},
		},
		{
			name:  "multi-bit codes spanning two bytes",
			text:  "ab",
			codes: map[rune]string{'a': "000", 'b': "111"},
			// "000111" padded to "00011100"
			expectedBits: []byte{0b00011100},
		},
		{
			name:  "codes that produce exactly two bytes",
			text:  "aabb",
			codes: map[rune]string{'a': "00", 'b': "11"},
			// "00001111"
			expectedBits: []byte{0b00001111},
		},
		{
			name:  "codes that produce exactly two bytes",
			text:  "aabbccd",
			codes: map[rune]string{'a': "00", 'b': "11", 'c': "10", 'd': "01"},
			// "00001111 10100101"
			expectedBits: []byte{0b00001111, 0b10100100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := writeBits(tt.text, tt.codes)

			require.Equal(t, len(tt.expectedBits), len(got), "writeBits() did not return same lengths")
			assert.Equal(t, tt.expectedBits, got, "writeBits not returned same slice of bits")
		})
	}
}

func TestGenerateCodes(t *testing.T) {
	tests := []struct {
		name          string
		root          *TreeNode
		expectedCodes map[rune]string
	}{
		{
			name:          "nil root produces no codes",
			root:          nil,
			expectedCodes: map[rune]string{},
		},
		{
			name:          "single leaf node",
			root:          &TreeNode{Char: 'a', Freq: 1},
			expectedCodes: map[rune]string{'a': ""},
		},
		{
			name: "two-leaf tree",
			root: &TreeNode{
				Freq:  3,
				Left:  &TreeNode{Char: 'a', Freq: 1},
				Right: &TreeNode{Char: 'b', Freq: 2},
			},
			expectedCodes: map[rune]string{'a': "0", 'b': "1"},
		},
		{
			name: "three-character tree",
			// Expected tree for 'a':1, 'b':2, 'c':3 (freq-ordered):
			//       (6)
			//      /   \
			//    (3)    c:3
			//   /   \
			//  a:1   b:2
			root: &TreeNode{
				Freq: 6,
				Left: &TreeNode{
					Freq:  3,
					Left:  &TreeNode{Char: 'a', Freq: 1},
					Right: &TreeNode{Char: 'b', Freq: 2},
				},
				Right: &TreeNode{Char: 'c', Freq: 3},
			},
			expectedCodes: map[rune]string{'a': "00", 'b': "01", 'c': "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codes := make(map[rune]string)
			generateCodes(tt.root, "", codes)

			assert.Equal(t, tt.expectedCodes, codes, "generateCodes did not return expected codes")
		})
	}
}
