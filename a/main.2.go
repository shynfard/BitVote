// Welcome to the gnark playground!
package main

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/uints"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	privateKey eddsa.PrivateKey  `gnark:",public"`
	PublicKey  eddsa.PublicKey   `gnark:",public"`
	Signature  eddsa.Signature   `gnark:",public"`
	Message    frontend.Variable `gnark:",public"`
}

func SHA3(api frontend.API, data []frontend.Variable) []frontend.Variable {
	hash, err := sha3.New512(api)
	if err != nil {
		panic(err)
	}
	for _, d := range data {
		bytes := api.ToBinary(d, 256)
		tb := make([]uints.U8, len(bytes))
		for i, b := range bytes {
			tb[i] = uints.U8{
				Val: b,
			}
		}
		hash.Write(tb)
	}
	sum := hash.Sum()
	return api.ToBinary(sum, 256)
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
// This circuit verifies an EdDSA signature.
func (circuit *Circuit) Define(api frontend.API) error {

	h, err := sha3.New512(api)
	if err != nil {
		return err
	}

	x := []uints.U8
	for i := 0; i < 32; i++ {
		h.Write([]uint8(circuit.Message))
	}
	h.Write([]uint8(circuit.privateKey.Bytes()))

	// pk, err := eddsa.GenerateKey()

	return nil
}
