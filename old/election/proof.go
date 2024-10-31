package election

import (
	"bytes"
	"fmt"
	"io"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/selector"
)

type Circuit struct {
	Secret frontend.Variable   `gnark:"x"` // the private input (secret)
	List   []frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	members := selector.KeyDecoder(api, circuit.Secret, circuit.List)
	temp := frontend.Variable(0)
	for i := 0; i < len(members); i++ {
		temp = api.Add(temp, members[i])
	}
	api.AssertIsEqual(temp, frontend.Variable(1))
	return nil
}

func main() {
	// login()

	circuit := Circuit{List: []frontend.Variable{frontend.Variable(3), frontend.Variable(2), frontend.Variable(4), frontend.Variable(5)}}

	// Compile the circuit
	ccs, err := frontend.Compile(ecc.BLS12_377.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println("Circuit compilation error:", err)
		return
	}

	fmt.Println("Circuit compiled successfully", circuit)

	// Create the Prover and Verifier keys, then generate the proof and verify

	// Setup the proving and verifying keys (trusted setup)
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}
	fmt.Println("Setup successful")

	// witness definition
	assignment := Circuit{Secret: frontend.Variable(5), List: []frontend.Variable{frontend.Variable(3), frontend.Variable(2), frontend.Variable(4), frontend.Variable(5)}}
	witness1, _ := frontend.NewWitness(&assignment, ecc.BLS12_377.ScalarField())
	publicWitness, _ := witness1.Public()
	m, _ := publicWitness.MarshalBinary()
	fmt.Println("Witness:", string(m))

	// Generate a proof
	proof, err := groth16.Prove(ccs, pk, witness1)
	if err != nil {
		fmt.Println("Proof generation error:", err)
		return
	}

	var proofBuf bytes.Buffer
	var proofWriter io.Writer = &proofBuf
	proof.WriteTo(proofWriter)

	var vkBuf bytes.Buffer
	var vkWriter io.Writer = &vkBuf
	vk.WriteTo(vkWriter)

	publicWitnessBuff, err := publicWitness.MarshalBinary()
	if err != nil {
		fmt.Println("Marshalling error:", err)
		return
	}

	// ------------------------------
	// ------------------------------
	// ------------------------------
	// ------------------------------
	// ------------------------------

	newProof := groth16.NewProof(ecc.BLS12_377)
	newProof.ReadFrom(&proofBuf)

	newVk := groth16.NewVerifyingKey(ecc.BLS12_377)
	newVk.ReadFrom(&vkBuf)

	newPublicWitness, _ := witness.New(ecc.BLS12_377.ScalarField()) //
	newPublicWitness.UnmarshalBinary(publicWitnessBuff)

	// Verify the proof
	err = groth16.Verify(newProof, newVk, newPublicWitness)
	if err != nil {
		fmt.Println("Proof verification failed:", err)
	} else {
		fmt.Println("Proof verification succeeded!")
	}

	// // Binary marshalling
	// data, err := witness.MarshalBinary()
	// if err != nil {
	// 	fmt.Println("Marshalling error:", err)
	// 	return
	// }
	// fmt.Println("Marshalled data:", len(data), string(data))
	// // Binary marshalling
	// data2, err := publicWitness.MarshalBinary()
	// if err != nil {
	// 	fmt.Println("Marshalling error:", err)
	// 	return
	// }
	// fmt.Println("Marshalled data:", len(data2), string(data2))
	// // Binary marshalling

}
