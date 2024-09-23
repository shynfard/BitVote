// Welcome to the gnark playground!
package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/uints"
)

type Circuit struct {
	In       []uints.U8
	Expected []uints.U8
}

func (c *Circuit) Define(api frontend.API) error {
	// h, err := sha3.New256(api)
	// if err != nil {
	// 	return err
	// }
	// uapi, err := uints.New[uints.U64](api)
	// if err != nil {
	// 	return err
	// }

	// h.Write(c.In)
	// res := h.Sum()
	// fmt.Println("res", res)

	// for i := range c.Expected {
	// 	uapi.ByteAssertEq(c.Expected[i], res[i])
	// }
	return nil

}

// func main() {

// 	in := make([]byte, 310)
// 	_, err := rand.Reader.Read(in)
// 	if err != nil {
// 		panic(err)
// 	}

// 	h := sha3256.New256()
// 	h.Write(in)
// 	expected := h.Sum(nil)
// 	assignment := &Circuit{
// 		In:       make([]uints.U8, len(in)),
// 		Expected: make([]uints.U8, len(expected)),
// 	}

// 	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})
// 	field := ecc.BN254.ScalarField()

// 	fmt.Printf("inner Ccs nbConstraints:%v, nbSecretWitness:%v, nbPublicInstance:%v\n", ccs.GetNbConstraints(), ccs.GetNbSecretVariables(), ccs.GetNbPublicVariables())

// 	// NB! UNSAFE! Use MPC.
// 	srs, srsLagrange, err := unsafekzg.NewSRS(ccs)
// 	if err != nil {
// 		panic(err)
// 	}
// 	pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
// 	if err != nil {
// 		panic(err)
// 	}
// 	witness, err := frontend.NewWitness(assignment, field)

// 	proof, err := plonk.Prove(ccs, pk, witness)
// 	if err != nil {
// 		fmt.Printf("1err:%v\n", err)
// 		panic(err)
// 	}
// 	publicWitness, err := witness.Public()
// 	if err != nil {
// 		fmt.Printf("2err:%v\n", err)
// 		panic(err)
// 	}
// 	err = plonk.Verify(proof, vk, publicWitness)
// 	if err != nil {
// 		fmt.Printf("3err:%v\n", err)
// 		panic(err)
// 	}

// }

func maian() {

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
