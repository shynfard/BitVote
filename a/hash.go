package main

import (
	"fmt"

	sha3256 "golang.org/x/crypto/sha3"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/uints"
)

// MembershipCircuit defines the circuit for membership proof
type MembershipCircuit struct {
	HiddenValue   []byte
	PublicSetHash []frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
func (circuit *MembershipCircuit) Define(api frontend.API) error {
	hashedHiddenValue := hashSet(circuit.HiddenValue, api)
	api.Println("PublicSetHash", circuit.PublicSetHash)
	api.Println("hashedHiddenValue", hashedHiddenValue)
	for i := range circuit.PublicSetHash {
		api.AssertIsEqual(hashedHiddenValue[i], circuit.PublicSetHash[i])
	}

	return nil
}

// Helper function to hash a set of values using MiMC
func hashSet(values []byte, api frontend.API) []frontend.Variable {
	hasher, err := sha3.New256(api)
	if err != nil {
		fmt.Println("Error creating sha3 256 hasher:", err)
		panic(err)
	}
	data := uints.NewU8Array(values)
	hasher.Write(data)
	hash := hasher.Sum()

	v := make([]frontend.Variable, len(hash))
	for i := range v {
		v[i] = hash[i].Val
	}
	return v
}

func main() {

	// in := make([]byte, 310)
	in := []byte{1, 2}
	fmt.Printf("in: %v\n", in)

	h := sha3256.New256()
	h.Write(in)
	expected := h.Sum(nil)
	expectedVar := make([]frontend.Variable, len(expected))
	for i := range expected {
		expectedVar[i] = frontend.Variable(expected[i])
	}
	fmt.Printf("expectedVar: %v\n", expectedVar)

	assignment := &MembershipCircuit{
		HiddenValue:   in,
		PublicSetHash: expectedVar,
	}
	witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	circuit := MembershipCircuit{
		HiddenValue:   in,
		PublicSetHash: make([]frontend.Variable, len(expected)),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	// r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, assignment)
	if err != nil {
		return
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return
	}
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Printf("verification failed\n")
		return
	}
	fmt.Printf("verification succeded\n")
}
