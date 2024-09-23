// Welcome to the gnark playground!
package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/uints"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	R frontend.Variable
	S frontend.Variable `gnark:",public"`
}

func SHA3(api frontend.API, data frontend.Variable) []frontend.Variable {
	hash, err := sha3.New512(api)
	if err != nil {
		panic(err)
	}
	dataBytes := api.ToBinary(data, 256)
	tb := make([]uints.U8, len(dataBytes))
	tb = uints.NewU8Array(dataBytes)
	// // uints.NewU8Array(in),
	tb := make([]uints.U8, len(bytes))
	for i, b := range bytes {
		tb[i] = uints.U8{
			Val: b,
		}
	}
	hash.Write(tb)
	sum := hash.Sum()
	fmt.Println("sum", sum)

	return api.ToBinary(sum, 256)
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
// This circuit verifies an EdDSA signature.
func (circuit *Circuit) Define(api frontend.API) error {

	h := SHA3(api, circuit.R)

	api.AssertIsEqual(h, circuit.S)

	return nil
}

func main() {
	// compiles the circuit into a R1CS
	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := Circuit{R: 2, S: 0}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("invalid proof")
	} else {
		println("valid proof")
	}
}
