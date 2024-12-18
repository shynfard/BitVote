package election

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"

	bn254Teddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"github.com/shynfard/BitVote/wallet"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"

	// sigeddsa "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type NIZKProof struct {
	c, z *big.Int
}

type Vote struct {
	poll *Poll

	vote          []byte
	encryptedVote [][]byte

	wallet     *wallet.Wallet
	privateKey *bn254Teddsa.PrivateKey

	publicWitnessBuff []byte
	vkBuf             []byte
	proofBuf          []byte

	signature []byte

	rand *big.Int

	publicKey []byte
}

type VotePublic struct {
	PollID            []byte   `json:"pollID"`
	EncryptedVote     [][]byte `json:"encryptedVote"`
	PublicWitnessBuff []byte   `json:"publicWitnessBuff"`
	VkBuf             []byte   `json:"vkBuf"`
	ProofBuf          []byte   `json:"proofBuf"`
	Signature         []byte   `json:"signature"`
	PublicKey         []byte   `json:"publicKey"`
}

// - set poll ID
// - set vote
// - encrypt vote with public key of poll creator
// - create a one-time key pair
// - create proof of authenticity (that one-time key pair is generated by master keys)
// - create proof of authorization (that master public key is in participants list)
// - calculate key image
// - sign vote with private spend key
func CreateVote(wallet wallet.Wallet, pollData []byte, vote []byte) *Vote {
	v := &Vote{}
	v.wallet = &wallet
	poll, err := LoadPoll(pollData)
	if err != nil {
		panic(err)
	}
	v.poll = poll
	v.vote = vote
	v.rand = randomBigInt(v.poll.homomorphicPublicKey.N)

	v.calculateEncryptedVote()

	v.privateKey = wallet.GenerateOneTimePair(v.rand.Bytes())

	v.calculateProof()

	return v

}

func (v *Vote) calculateEncryptedVote() {
	for _, dataVote := range v.vote {
		enc, err := paillier.EncryptWithNonce(v.poll.homomorphicPublicKey, v.rand, []byte{dataVote})
		if err != nil {
			panic(err)
		}
		v.encryptedVote = append(v.encryptedVote, enc.Bytes())
	}
}

type Circuit struct {
	Random        frontend.Variable `gnark:",public"`
	PublicKey     eddsa.PublicKey
	ListPublicKey []eddsa.PublicKey `gnark:",public"`
	Signature     eddsa.Signature   `gnark:",public"`
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

func (v *Vote) calculateProof() {
	circuit := Circuit{
		ListPublicKey: make([]eddsa.PublicKey, len(v.poll.participants)),
	}

	// Compile the circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println("Circuit compilation error:", err)
		return
	}

	// Setup the proving and verifying keys (trusted setup)
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}

	// witness definition
	msg := v.privateKey.PublicKey.Bytes()[:31]
	// msg := []byte{0xde, 0xad, 0xf0, 0x0d}

	signature, err := v.wallet.Sign(msg)
	if err != nil {
		fmt.Println("Error signing message -- :", err)
		panic(err)
	}

	// declare the witness
	assignment := Circuit{}

	// assign message value
	assignment.Random = msg

	// assign public key values
	assignment.PublicKey.Assign(1, v.wallet.GetPublicKey().Bytes()[:32])

	// assign signature values
	assignment.Signature.Assign(1, signature)

	ListPublicKey := make([]eddsa.PublicKey, len(v.poll.participants))
	for i, participant := range v.poll.participants {
		ListPublicKey[i].Assign(1, participant)
	}
	assignment.ListPublicKey = ListPublicKey

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

	v.publicWitnessBuff = publicWitnessBuff
	v.vkBuf = vkBuf.Bytes()
	v.proofBuf = proofBuf.Bytes()

}

func (v *Vote) GetHash() []byte {
	h := sha256.New()
	h.Write(v.poll.pollID)
	for _, encVote := range v.encryptedVote {
		h.Write(encVote)
	}
	h.Write(v.privateKey.PublicKey.Bytes())
	h.Write(v.publicWitnessBuff)
	h.Write(v.vkBuf)
	h.Write(v.proofBuf)
	h.Write(v.signature)
	h.Write(v.rand.Bytes())
	return h.Sum(nil)
}

func (v *Vote) GetVote() []byte {
	data := map[string]interface{}{
		"pollId":            v.poll.pollID,
		"encryptedVote":     v.encryptedVote,
		"publicWitnessBuff": v.publicWitnessBuff,
		"vkBuf":             v.vkBuf,
		"proofBuf":          v.proofBuf,
		"singature":         v.signature,
		"publicKey":         v.privateKey.PublicKey.Bytes(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return jsonData
}
func LoadVote(data []byte, p Poll) (*Vote, error) {
	var v VotePublic
	err := json.Unmarshal(data, &v)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling PollPublic: %w", err)
	}
	vote := Vote{
		publicWitnessBuff: v.PublicWitnessBuff,
		vkBuf:             v.VkBuf,
		proofBuf:          v.ProofBuf,
		signature:         v.Signature,
		publicKey:         v.PublicKey,
		poll:              &p,
	}
	return &vote, nil
}

func randomBigInt(max *big.Int) *big.Int {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	return r
}

func (v *Vote) VerifyVote() bool {
	newProof := groth16.NewProof(ecc.BN254)
	newProof.ReadFrom(bytes.NewReader(v.proofBuf))

	newVk := groth16.NewVerifyingKey(ecc.BN254)
	newVk.ReadFrom(bytes.NewReader(v.vkBuf))

	newPublicWitness, _ := witness.New(ecc.BN254.ScalarField()) //
	newPublicWitness.UnmarshalBinary(v.publicWitnessBuff)

	var assignment Circuit

	// assign message value
	assignment.Random = v.publicKey[:31]
	assignment.Signature.Assign(1, v.signature)
	ListPublicKey := make([]eddsa.PublicKey, len(v.poll.participants))
	for i, participant := range v.poll.participants {
		ListPublicKey[i].Assign(1, participant)
	}
	assignment.ListPublicKey = ListPublicKey
	witness1, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness1.Public()

	// Verify the proof
	err := groth16.Verify(newProof, newVk, publicWitness)
	if err != nil {
		fmt.Println("Proof verification failed:", err)
	} else {
		fmt.Println("Proof verification succeeded!")
	}
	return true
}
