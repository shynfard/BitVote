package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bn254T "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	bn254Teddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
	"github.com/shynfard/BitVote/wallet"
	"golang.org/x/crypto/blake2b"
)

const (
	sizeFr         = fr.Bytes
	sizePublicKey  = sizeFr
	sizeSignature  = 2 * sizeFr
	sizePrivateKey = 2*sizeFr + 32
)

// Login or create new account
func login() (*wallet.Wallet, error) {

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
			return nil, err
		}
		names := string(fileData[:n])
		wallet.Load(names)
	}

	return &wallet, nil

}

// PrivateKey private key of an eddsa instance
type PrivateKey struct {
	PublicKey bn254Teddsa.PublicKey // copy of the associated public key
	scalar    [sizeFr]byte          // secret scalar, in big Endian
	randSrc   [32]byte              // source
}

type Circuit struct {
	Random        frontend.Variable `gnark:",public"`
	PublicKey     eddsa.PublicKey
	ListPublicKey []eddsa.PublicKey `gnark:",public"`
	Signature     eddsa.Signature   `gnark:",public"`
}

// GenerateKey generates a public and private key pair.
func GenerateKey(data [32]byte) (*PrivateKey, error) {
	c := bn254T.GetEdwardsCurve()

	var pub bn254Teddsa.PublicKey
	var priv PrivateKey
	// hash(h) = private_key || random_source, on 32 bytes each
	h := blake2b.Sum512(data[:])
	for i := 0; i < 32; i++ {
		priv.randSrc[i] = h[i+32]
	}

	// prune the key
	// https://tools.ietf.org/html/rfc8032#section-5.1.5, key generation
	h[0] &= 0xF8
	h[31] &= 0x7F
	h[31] |= 0x40

	// reverse first bytes because setBytes interpret stream as big endian
	// but in eddsa specs s is the first 32 bytes in little endian
	for i, j := 0, sizeFr-1; i < sizeFr; i, j = i+1, j-1 {
		priv.scalar[i] = h[j]
	}

	var bScalar big.Int
	bScalar.SetBytes(priv.scalar[:])
	pub.A.ScalarMultiplication(&c.Base, &bScalar)

	priv.PublicKey = pub

	return &priv, nil
}

func (circuit *Circuit) Define(api frontend.API) error {

	curve, err := twistededwards.NewEdCurve(api, 1)
	if err != nil {
		return err
	}

	hash, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	eddsa.Verify(curve, circuit.Signature, circuit.Random, circuit.PublicKey, &hash)

	temp := frontend.Variable(0)
	for i := 0; i < len(circuit.ListPublicKey); i++ {
		temp = api.Add(temp, api.Select(
			api.And(
				api.IsZero(api.Cmp(circuit.ListPublicKey[i].A.X, circuit.PublicKey.A.X)),
				api.IsZero(api.Cmp(circuit.ListPublicKey[i].A.Y, circuit.PublicKey.A.Y)),
			), frontend.Variable(1), frontend.Variable(0)))
	}
	api.AssertIsEqual(temp, frontend.Variable(1))
	return nil
}

func main() {
	wallet, _ := login()

	circuit := Circuit{
		ListPublicKey: make([]eddsa.PublicKey, 3),
	}

	// Compile the circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println("Circuit compilation error:", err)
		return
	}

	fmt.Println("Circuit compiled successfully", ccs)

	// Create the Prover and Verifier keys, then generate the proof and verify

	// Setup the proving and verifying keys (trusted setup)
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}
	fmt.Println("Setup successful")

	// witness definition
	msg := []byte{0xde, 0xad, 0xf0, 0x0d}

	signature, err := wallet.Sign(msg)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	// declare the witness
	var assignment Circuit

	// assign message value
	assignment.Random = msg

	// assign public key values
	assignment.PublicKey.Assign(1, wallet.GetPublicKey().Bytes()[:32])

	// assign signature values
	assignment.Signature.Assign(1, signature)

	var publicKey1 eddsa.PublicKey
	s1 := make([]byte, 32)
	rand.Read(s1)
	publicKey1.Assign(1, s1)
	var publicKey2 eddsa.PublicKey
	s2 := make([]byte, 32)
	rand.Read(s2)
	publicKey2.Assign(1, s2)
	var publicKey3 eddsa.PublicKey
	s3 := make([]byte, 32)
	rand.Read(s3)
	publicKey3.Assign(1, s3)

	assignment.ListPublicKey = []eddsa.PublicKey{publicKey1, publicKey2, assignment.PublicKey}
	fmt.Println(assignment)

	witness1, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
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

	// // ------------------------------
	// // ------------------------------
	// // ------------------------------
	// // ------------------------------
	// // ------------------------------

	newProof := groth16.NewProof(ecc.BN254)
	newProof.ReadFrom(&proofBuf)

	newVk := groth16.NewVerifyingKey(ecc.BN254)
	newVk.ReadFrom(&vkBuf)

	newPublicWitness, _ := witness.New(ecc.BN254.ScalarField()) //
	newPublicWitness.UnmarshalBinary(publicWitnessBuff)

	// Verify the proof
	err = groth16.Verify(newProof, newVk, publicWitness)
	if err != nil {
		fmt.Println("Proof verification failed:", err)
	} else {
		fmt.Println("Proof verification succeeded!")
	}

}
