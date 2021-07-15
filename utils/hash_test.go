package utils

import (
	"testing"
)

func TestSha256Hex(t *testing.T) {
	type sha256Test struct {
		input  string
		output string
	}
	testCases := []sha256Test{
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}
	for _, testCase := range testCases {
		result := Sha256Hex(testCase.input)
		if result != testCase.output {
			t.Errorf("incorrect SHA-256 hash of '%s' => '%s' (expected %s)\n", testCase.input, result, testCase.output)
		}
	}
}
