// Welcome to the gnark playground!
package main

import (
	"encoding/hex"
	"fmt"

	sha3256 "golang.org/x/crypto/sha3"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/test/unsafekzg"
)

const MaxLength = 128

type Circuit struct {
	Input       [MaxLength]uints.U8
	InputLength frontend.Variable     //indicates actual number of bytes
	Output      [32]frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	var digests [32]frontend.Variable

	h, err := sha3.New256(api)
	if err != nil {
		return err
	}
	h.Write(circuit.Input[:])
	digestU8 := h.Sum()
	for i := range circuit.Output {
		api.AssertIsEqual(circuit.Output[i], digestU8[i].Val)
		digests[i] = digestU8[i].Val
	}

	return nil

}

func main() {
	field := ecc.BN254.ScalarField()

	// compiles the circuit into a R1CS

	msg := []byte("hello, world1234")

	digest := sha3256.Sum256(msg)
	fmt.Printf("%v\n", hex.EncodeToString(digest[:]))

	var input [MaxLength]uints.U8
	var output [32]frontend.Variable
	for i := 0; i < len(msg); i++ {
		input[i] = uints.NewU8(msg[i])
	}
	for i, d := range digest {
		output[i] = d
	}

	assignment := &Circuit{
		Input:       input,
		InputLength: len(msg),
		Output:      output,
	}
	fmt.Printf("nbInput:%v, nbOutput:%v\n", len(assignment.Input), len(assignment.Output))
	ccs, err := frontend.Compile(field, scs.NewBuilder, &Circuit{})

	if err != nil {
		panic(err)
	}
	fmt.Printf("inner Ccs nbConstraints:%v, nbSecretWitness:%v, nbPublicInstance:%v\n", ccs.GetNbConstraints(), ccs.GetNbSecretVariables(), ccs.GetNbPublicVariables())

	// NB! UNSAFE! Use MPC.
	srs, srsLagrange, err := unsafekzg.NewSRS(ccs)
	if err != nil {
		panic(err)
	}
	pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
	if err != nil {
		panic(err)
	}
	witness, err := frontend.NewWitness(assignment, field)

	proof, err := plonk.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Printf("1err:%v\n", err)
		panic(err)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		fmt.Printf("2err:%v\n", err)
		panic(err)
	}
	err = plonk.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Printf("3err:%v\n", err)
		panic(err)
	}
}
