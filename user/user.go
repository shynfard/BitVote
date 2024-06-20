package user

import (
	"crypto/sha512"
	"fmt"

	"filippo.io/edwards25519"
)

// One-time key pair
type OneTimePair struct {
	private *edwards25519.Scalar
	public  *edwards25519.Point
}

type User struct {
	privateViewKey  *edwards25519.Scalar
	privateSpendKey *edwards25519.Scalar
	publicViewKey   *edwards25519.Point
	publicSpendKey  *edwards25519.Point

	// last One-time key
	oneTinePairs []*OneTimePair
}

func NewUser() (*User, error) {
	privateViewKey, err := generatePrivateKey()
	if err != nil {
		fmt.Println("Error generating private view key:", err)
		return nil, err
	}
	privateSpendKey, err := generatePrivateKey()
	if err != nil {
		fmt.Println("Error generating private spend key:", err)
		return nil, err
	}

	// Derive public view key and public spend key
	publicViewKey := derivePublicKey(privateViewKey)
	publicSpendKey := derivePublicKey(privateSpendKey)

	u := &User{
		privateViewKey:  privateViewKey,
		privateSpendKey: privateSpendKey,
		publicViewKey:   publicViewKey,
		publicSpendKey:  publicSpendKey,
	}
	return u, nil
}

// Generate one-time public key from private view key and public spend key
func (u *User) GenerateOneTimeKeyPair(index uint64) error {
	// Convert the index to a 32-byte array (little-endian)
	indexBytes := [32]byte{}
	for i := 0; i < 8; i++ {
		indexBytes[i] = byte((index >> (8 * i)) & 0xff)
	}

	// Calculate the hash of the private view key and index
	h := sha512.New()
	h.Write(u.privateViewKey.Bytes())
	h.Write(indexBytes[:])
	hash := h.Sum(nil)

	// Reduce the hash mod l to obtain a scalar (one-time private key)
	var oneTimePrivateKey edwards25519.Scalar
	oneTimePrivateKey.SetBytesWithClamping(hash[:32])

	// Calculate the one-time public key: oneTimePrivateKey * G + publicSpendKey
	var oneTimePublicKey edwards25519.Point
	oneTimePublicKey.ScalarBaseMult(&oneTimePrivateKey)
	oneTimePublicKey.Add(&oneTimePublicKey, u.publicSpendKey)

	otp := &OneTimePair{
		private: &oneTimePrivateKey,
		public:  &oneTimePublicKey,
	}
	u.oneTinePairs = append(u.oneTinePairs, otp)
	return nil
}

func (u *User) GetLastOneTimePublicKey() *edwards25519.Point {
	return u.oneTinePairs[len(u.oneTinePairs)-1].public
}
func (u *User) getLastOneTimePair() *OneTimePair {
	return u.oneTinePairs[len(u.oneTinePairs)-1]
}

func (u *User) Sign(message []byte) []byte {
	// Calculate the hash of the message
	h := sha512.New()
	h.Write(message)
	hash := h.Sum(nil)

	lastPair := u.getLastOneTimePair()
	// Calculate the signature: oneTimePrivateKey + hash
	var signature edwards25519.Scalar
	signature.SetBytesWithClamping(hash[:32])
	signature.Add(lastPair.private, &signature)

	return signature.Bytes()
}

// Sign a message with a private key
func (u *User) signMessage(message []byte) ([]byte, error) {
	lastPair := u.getLastOneTimePair()

	// Hash the private key and the message
	h := sha512.New()
	h.Write(lastPair.private.Bytes())
	h.Write(message)
	hash := h.Sum(nil)

	// Reduce the hash mod l to obtain a scalar (signature)
	var signature edwards25519.Scalar
	signature.SetBytesWithClamping(hash[:32])

	return signature.Bytes(), nil
}

// Verify a signature with a public key
func verifySignature(publicKey *edwards25519.Point, message, signature []byte) bool {
	// Hash the public key, the message, and the signature
	h := sha512.New()
	h.Write(publicKey.Bytes())
	h.Write(message)
	h.Write(signature)
	hash := h.Sum(nil)

	// Reduce the hash mod l to obtain a scalar (computed signature)
	var computedSignature edwards25519.Scalar
	computedSignature.SetBytesWithClamping(hash[:32])

	// Reconstruct the public key from the signature and the base point
	var reconstructedPublicKey edwards25519.Point
	reconstructedPublicKey.ScalarBaseMult(&computedSignature)

	return reconstructedPublicKey.Equal(publicKey) == 1
}
