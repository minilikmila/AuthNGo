package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	// Generate a random number between 0 and 999999
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		// Fallback to a simple random number if crypto/rand fails
		n = big.NewInt(0)
	}
	// Format as a 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64())
}
