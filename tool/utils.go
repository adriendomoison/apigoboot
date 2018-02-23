/*
	msutils add string generation, modification and conversion capabilities
*/
package tool

import (
	"time"
	"math/rand"
	"regexp"
	"strings"
)

// ToSnakeCase change a string to it's snake case version
func ToSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

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