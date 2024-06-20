package user

import (
	"crypto/rand"

	"filippo.io/edwards25519"
)

// Generate a new private key
func generatePrivateKey() (*edwards25519.Scalar, error) {
	privateKeyBytes := [32]byte{}
	_, err := rand.Read(privateKeyBytes[:])
	if err != nil {
		return nil, err
	}
	privateKey := new(edwards25519.Scalar)
	privateKey.SetBytesWithClamping(privateKeyBytes[:])
	return privateKey, nil
}

// Derive the public key from a private key
func derivePublicKey(privateKey *edwards25519.Scalar) *edwards25519.Point {
	publicKey := new(edwards25519.Point).ScalarBaseMult(privateKey)
	return publicKey
}
