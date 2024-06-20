package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type CubicCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X    frontend.Variable   `gnark:",secret"`
	List []frontend.Variable `gnark:"-"`
	Y    frontend.Variable   `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *CubicCircuit) Define(api frontend.API) error {
	// x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	// api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))

	fmt.Println("circuit y", circuit.Y)
	fmt.Println(circuit.List)
	for _, item := range circuit.List {
		isEqual := api.IsZero(api.Cmp(circuit.X, item))
		circuit.Y = api.Add(circuit.Y, isEqual)
	}
	fmt.Println(circuit.Y)
	api.AssertIsEqual(circuit.Y, frontend.Variable(1))
	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit CubicCircuit
	circuit = CubicCircuit{Y: 0, List: []frontend.Variable{1, 2, 3}}
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := CubicCircuit{X: 2, Y: 0, List: []frontend.Variable{1, 2, 3}}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)

	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("invalid proof")
	} else {
		println("valid proof")
	}

}
