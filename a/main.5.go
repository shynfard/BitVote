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

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier

// MAX_LENGTH is the maximum length of the input
const MaxLength = 128

type Circuit struct {
	Input       []uints.U8
	InputLength frontend.Variable //indicates actual number of bytes
	Output      []uints.U8        `gnark:",public"`
}

func SHA3(api frontend.API, data []uint8) []frontend.Variable {

	hash, err := sha3.New512(api)
	if err != nil {
		panic(err)
	}
	tb := make([]uints.U8, len(data))
	tb = uints.NewU8Array(data)
	hash.Write(tb)
	sum := hash.Sum()
	fmt.Println("sum", sum)

	return api.ToBinary(sum, 256)
}

// Define declares the circuit logicircuit. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
// This circuit verifies an EdDSA signature.
func (circuit *Circuit) Define(api frontend.API) error {
	var digests []uints.U8

	//calculated all possible digest
	// for i := 0; i < MaxLength; i++ {
	// 	h, err := sha3.New256(api)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	h.Write(circuit.Input[:i])
	// 	o := h.Sum()
	// 	digests[i] = o
	// }

	h, err := sha3.New256(api)
	if err != nil {
		return err
	}

	h.Write(circuit.Input)
	o := h.Sum()
	digests = o

	// //select
	// var res [64]frontend.Variable
	// for i := 0; i < len(circuit.Output); i++ {
	// 	var muxs []frontend.Variable
	// 	for j := 0; j < MaxLength; j++ {
	// 		muxs = append(muxs, digests[j][i].Val)
	// 	}
	// 	res[i] = selector.Mux(api, circuit.InputLength, muxs...)
	// }

	for i := range digests {
		api.AssertIsEqual(digests[i], circuit.Output[i])
	}
	return nil

}

func main() {
	field := ecc.BN254.ScalarField()

	// compiles the circuit into a R1CS

	msg := []byte("hello, world")

	digest := sha3256.Sum256(msg)
	fmt.Printf("%v\n", hex.EncodeToString(digest[:]))

	var input []uints.U8
	var output []uints.U8
	for i := 0; i < min(MaxLength, len(msg)); i++ {
		input[i] = uints.NewU8(msg[i])
	}
	for i, d := range digest {
		output[i] = uints.NewU8(d)
	}

	assignment := &Circuit{
		Input:       input,
		InputLength: len(msg),
		Output:      output,
	}
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
