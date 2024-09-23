package wallet

import (
	"filippo.io/edwards25519"
)

// One-time key pair
type OneTimePair struct {
	private *edwards25519.Scalar
	public  *edwards25519.Point
}

type Wallet struct {
	privateViewKey  *edwards25519.Scalar
	privateSpendKey *edwards25519.Scalar
	publicViewKey   *edwards25519.Point
	publicSpendKey  *edwards25519.Point

	// list used One-time key
	oneTimePairs []*OneTimePair
}

// Generate a new wallet
func (w *Wallet) Generate() {
	// generate private view key
	privateViewKey, _ := generatePrivateKey()
	w.privateViewKey = privateViewKey
	w.publicViewKey = derivePublicKey(privateViewKey)

	// generate private spend key
	privateSpendKey, _ := generatePrivateKey()
	w.privateSpendKey = privateSpendKey
	w.publicSpendKey = derivePublicKey(privateSpendKey)
}

// load wallet
func (w *Wallet) Load(keys []byte) {
	w.privateViewKey = edwards25519.NewScalar()
	w.privateViewKey.SetUniformBytes(keys[:64])

	w.privateSpendKey = edwards25519.NewScalar()
	w.privateSpendKey.SetUniformBytes(keys[64:])

	w.publicSpendKey = derivePublicKey(w.privateSpendKey)
	w.publicViewKey = derivePublicKey(w.privateViewKey)
}

// get public key
func (w *Wallet) GetPublicKey() []byte {
	return w.publicSpendKey.Bytes()
}

// get public key
func (w *Wallet) GetPrivateKey() []byte {
	return w.privateSpendKey.Bytes()
}

// generate one-time key pair
func (w *Wallet) GenerateOneTimePair(rand []byte) (private *edwards25519.Scalar, public *edwards25519.Point) {
	// generate private key
	privateKey, _ := generateOneTimePrivateKey(w.privateSpendKey, rand)
	// generate public key
	publicKey := derivePublicKey(privateKey)

	// append to list
	w.oneTimePairs = append(w.oneTimePairs, &OneTimePair{private: privateKey, public: publicKey})

	return privateKey, publicKey
}
