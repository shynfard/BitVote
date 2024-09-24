package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/selector"
	"github.com/shynfard/BitVote/wallet"
)

// Login or create new account
func login() {

	wallet := wallet.Wallet{}

	// open file to read
	file, err := os.Open("wallet.txt")
	if err != nil {
		fmt.Println("No wallet.txt found, creating new wallet")
		names := wallet.Generate()
		file, err = os.Create("wallet.txt")
		if err != nil {
			fmt.Println("Error creating wallet.txt")
		}
		fmt.Println("Save:", names)
		fmt.Println("Len:", len(names))
		_, err = file.Write(bytes.NewBufferString(names).Bytes())
		if err != nil {
			fmt.Println("Error saving wallet.txt")
		}
	} else {
		// read file
		fileData := make([]byte, 2045)
		n, err := file.Read(fileData)
		if err != nil {
			fmt.Println("Error reading wallet.txt")
			return
		}
		names := string(fileData[:n])
		wallet.Load(names)
	}

}

type Circuit struct {
	Secret frontend.Variable   `gnark:",secret"` // the private input (secret)
	List   []frontend.Variable // the public list
}

func (circuit *Circuit) Define(api frontend.API) error {
	// Assume list contains 4 elements for simplicity, but it can be extended
	isMember := selector.IsInSlice(api, circuit.Secret, circuit.List)

	// Assert that isMember == 1 (i.e., the secret is in the list)
	api.AssertIsEqual(isMember, 1)

	return nil
}

func main() {
	// login()

	var circuit Circuit

	// Define the list of public inputs (public list)
	circuit.List = make([]frontend.Variable, 4)
	circuit.List[0] = 2
	circuit.List[1] = 3
	circuit.List[2] = 5
	circuit.List[3] = 7

	// Compile the circuit
	r1cs, err := frontend.Compile(ecc.BLS12_377.BaseField(), backend.GROTH16, &circuit)
	if err != nil {
		fmt.Println("Circuit compilation error:", err)
		return
	}

	fmt.Println("Circuit compiled successfully")

	// Create the Prover and Verifier keys, then generate the proof and verify

	// Setup the proving and verifying keys (trusted setup)
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}

	// Witness: set the secret to a value that is in the public list (e.g., 3)
	var witness Circuit
	witness.Secret = 3

	// Generate a proof
	proof, err := groth16.Prove(r1cs, pk, &witness)
	if err != nil {
		fmt.Println("Proof generation error:", err)
		return
	}

	// Create a public input (the public list)
	publicWitness := Circuit{
		List: []frontend.Variable{2, 3, 5, 7}, // public list
	}

	// Verify the proof
	err = groth16.Verify(proof, vk, &publicWitness)
	if err != nil {
		fmt.Println("Proof verification failed:", err)
	} else {
		fmt.Println("Proof verification succeeded!")
	}

}
