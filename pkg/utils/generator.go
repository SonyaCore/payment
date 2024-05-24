package utils

import (
	"math/rand"
	"time"
)

// characterSet are being used for generating random string in GenerateDiscount func
var characterSet = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")

// GenerateDiscount generates a random code of the specified length.
func GenerateDiscount(length int) string {
	rand.NewSource(time.Now().UnixNano())
	result := make([]rune, length)
	for i := range result {
		result[i] = characterSet[rand.Intn(len(characterSet))]
	}
	return string(result)
}
