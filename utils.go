package democracy

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"math/rand"
	"time"
)

const (
	LEN = 10
)

func RandGenerator() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// Define the range
	min := 1000000
	max := 9999999

	// Generate a random integer between min and max
	return r.Intn(max-min+1) + min
}

func GenerateShortID() (string, error) {
	// Create a byte slice to hold the random bytes
	bytes := make([]byte, LEN)

	// Fill the byte slice with random bytes
	_, err := cryptoRand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a URL-safe base64 string
	shortID := base64.URLEncoding.EncodeToString(bytes)

	// Truncate the string to the desired length
	if len(shortID) > LEN {
		shortID = shortID[:LEN]
	}

	return shortID, nil
}

func Ceil(a, b int) int {
	if a%b == 0 {
		return a / b
	} else {
		return (a / b) + 1
	}
}

func Unmarshal() {

}
