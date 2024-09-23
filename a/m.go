package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

// MembershipCircuit defines the circuit for membership proof
type MembershipCircuit struct {
	HiddenValue   frontend.Variable
	PublicSetHash frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
func (circuit *MembershipCircuit) Define(api frontend.API) error {
	// Hash the hidden value
	hashedHiddenValue := hashSet([]frontend.Variable{circuit.HiddenValue}, api)
	// Enforce that the hash of the hidden value is equal to the public set hash
	api.AssertIsEqual(hashedHiddenValue, circuit.PublicSetHash)

	return nil
}

// Helper function to hash a set of values using MiMC
func hashSet(values []frontend.Variable, api frontend.API) frontend.Variable {
	hasher, err := mimc.NewMiMC(api)
	if err != nil {
		fmt.Println("Error creating MiMC hasher:", err)
		return frontend.Variable(0)
	}
	for _, v := range values {
		hasher.Write(v)
	}
	return hasher.Sum()
}

func main() {

	// Compile the circuit
	var circuit MembershipCircuit
	css, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println("Error compiling circuit:", err)
		return
	}

	// Example public set and hidden value
	publicSet := []frontend.Variable{
		frontend.Variable(1),
		frontend.Variable(2),
		frontend.Variable(3),
	}
	hiddenValue := frontend.Variable(2)

	// Hash the public set
	api := frontend.NewBuilder(r1cs.NewBuilder)
	if err != nil {
		fmt.Println("Error creating builder:", err)
		return
	}
	publicSetHash := hashSet(publicSet, builder)

	// Create a witness
	witness := MembershipCircuit{
		HiddenValue:   hiddenValue,
		PublicSetHash: publicSetHash,
	}

	// Generate proving and verifying keys
	pk, vk, err := groth16.Setup(css)
	if err != nil {
		fmt.Println("Error during setup:", err)
		return
	}

	// Generate proof
	proof, err := groth16.Prove(r1cs, pk, &witness)
	if err != nil {
		fmt.Println("Error generating proof:", err)
		return
	}

	// Verify proof
	isValid, err := groth16.Verify(proof, vk, &witness)
	if err != nil {
		fmt.Println("Error verifying proof:", err)
		return
	}

	if isValid {
		fmt.Println("Proof is valid!")
	} else {
		fmt.Println("Proof is invalid.")
	}
}
