package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

func GenerateTradeID() string {
	rand.Seed(time.Now().UnixNano())

	randomLetters := make([]byte, 2)
	for i := range randomLetters {
		randomLetters[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	randomNumbers := make([]byte, 10)
	for i := range randomNumbers {
		randomNumbers[i] = numberBytes[rand.Intn(len(numberBytes))]
	}

	return string(randomLetters) + string(randomNumbers)
}

func StringWithCharset(length int, charset string) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func Random8DigitNumber(r *rand.Rand) int {
	// Returns a random 8-digit number
	return r.Intn(90000000) + 10000000
}
