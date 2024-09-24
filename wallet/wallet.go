package wallet

import (
	"crypto/sha256"
	"math/rand"
	"strings"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/ecdsa"
	"github.com/wordgen/wordlists/names"
)

// One-time key pair
type OneTimePair struct {
	privateKey *ecdsa.PrivateKey
}

type Wallet struct {
	privateViewKey  *ecdsa.PrivateKey
	privateSpendKey *ecdsa.PrivateKey

	names string

	oneTimePairs []*OneTimePair
}

// Generate a new wallet
func (w *Wallet) Generate() string {

	s := ""
	for i := 0; i < 20; i++ {
		s += names.Mixed[rand.Intn(len(names.Mixed)-1)]
		if i < 19 {
			s += " "
		}
	}
	w.Load(s)
	return s
}

// load wallet
func (w *Wallet) Load(names string) {

	w.names = names
	// derive private key from names
	x := strings.Split(names, " ")

	s1 := ""
	for i := 0; i < 10; i++ {
		s1 += x[i]
	}
	s2 := ""
	for i := 0; i < 10; i++ {
		s2 += x[i+5]
	}

	seed1 := []byte(s1)
	seed2 := []byte(s2)
	hash1 := sha256.Sum256(seed1)
	hash2 := sha256.Sum256(seed2)
	rng1 := rand.New(rand.NewSource(int64(hash1[0])))
	rng2 := rand.New(rand.NewSource(int64(hash2[0])))

	privateKey, err := ecdsa.GenerateKey(rng1)
	if err != nil {
		return
	}
	w.privateViewKey = privateKey

	privateKey, err = ecdsa.GenerateKey(rng2)
	if err != nil {
		return
	}
	w.privateSpendKey = privateKey
}

// generate one-time key pair
func (w *Wallet) GenerateOneTimePair(randInput []byte) (key *ecdsa.PrivateKey) {

	h := sha256.New()
	h.Write(w.privateSpendKey.Bytes())
	h.Write(randInput)
	privateKeyBytes := h.Sum(nil)
	hash1 := sha256.Sum256(privateKeyBytes)
	rng1 := rand.New(rand.NewSource(int64(hash1[0])))

	privateKey, err := ecdsa.GenerateKey(rng1)
	if err != nil {
		return
	}
	// append to list
	w.oneTimePairs = append(w.oneTimePairs, &OneTimePair{privateKey: privateKey})
	return privateKey
}
