package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// source is the random source used to generate random values
var source = rand.NewSource(time.Now().Unix())

// RandHexStr generates a random string of the given length
func RandHexStr(length uint) string {
	builder := strings.Builder{}
	for uint(builder.Len()) < length {
		builder.WriteString(RandHexStrInt64())
	}
	str := builder.String()
	return str[:length]
}

// RandHexStrInt64 generates a random 64-bit integer, and returns the hexadecimal-encoded string representation. The
// length of the returned string will always be exactly 16 bytes. As needed, the returned string will be padded with
// leading zeros
func RandHexStrInt64() string {
	str := strconv.FormatInt(source.Int63(), 16)
	// All random strings should be 16 characters long. We need to add leading zeros to enforce this
	for len(str) < 16 {
		str = "0" + str
	}
	return str
}
