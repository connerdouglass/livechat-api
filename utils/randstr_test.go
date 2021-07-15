package utils

import (
	"testing"
)

func TestRandHexStr(t *testing.T) {

	// Generate a bunch of the same length
	for i := 0; i < 100; i++ {
		result := RandHexStrInt64()
		if len(result) != 16 {
			t.Errorf("random hex string length is inconsistent. should always be 16")
		}
	}

	// Generate strings of many different lengths
	for i := 0; i < 100; i++ {
		result := RandHexStr(uint(i))
		if len(result) != i {
			t.Errorf("random hex string of length: %d expected %d", len(result), i)
		}
	}

}
