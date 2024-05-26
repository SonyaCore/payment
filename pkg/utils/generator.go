package utils

import (
	"crypto/rand"
	"math/big"
)

// characterSet are being used for generating random string in GenerateDiscount func
var characterSet = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")

// GenerateDiscount generates a random code of the specified length using CSPRNG.
func GenerateDiscount(length int) (string, error) {
	result := make([]rune, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(characterSet))))
		if err != nil {
			return "", err
		}
		result[i] = characterSet[index.Int64()]
	}
	return string(result), nil
}
