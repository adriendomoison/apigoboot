/*
	msutils add string generation, modification and conversion capabilities
*/
package tool

import (
	"time"
	"math/rand"
)

// GenerateRandomString create a random string of the requested length using the hexadecimal symbols
func GenerateRandomString(strLen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEF0123456789"
	result := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
