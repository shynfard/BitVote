package election

import (
	"crypto/rand"
	"log"
	"math/big"
)

func randomBigInt(max *big.Int) *big.Int {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	return r
}
