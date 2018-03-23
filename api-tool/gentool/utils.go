// Package gentool add string generation, modification and conversion capabilities
package gentool

import (
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
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

// GenerateCleanString take a string with symboles and extra spaces and return an alpha numerical string
func GenerateCleanString(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	alphanumericalString := reg.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(alphanumericalString), " ")
}
