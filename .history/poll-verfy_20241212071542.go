// package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/ed25519"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"runtime"
// 	"time"

// 	ced "crypto/ed25519"
// 	rrand "crypto/rand"
// 	"crypto/sha256"
// 	"encoding/json"
// 	"strings"

// 	"github.com/consensys/gnark-crypto/ecc"
// 	"github.com/consensys/gnark/backend/groth16"
// 	"github.com/consensys/gnark/frontend"
// 	"github.com/consensys/gnark/frontend/cs/r1cs"
// 	"github.com/consensys/gnark/std/algebra/native/twistededwards"
// 	"github.com/consensys/gnark/std/hash/mimc"
// 	"github.com/consensys/gnark/std/signature/eddsa"
// 	paillier "github.com/roasbeef/go-go-gadget-paillier"
// 	"github.com/shynfard/BitVote/old/wallet"

// 	"bytes"
// 	"io"
// 	"math/big"
// )

// // Poll represents an election poll.
// type Poll struct {
// 	Creator               wallet.Wallet
// 	creatorPublicKey      []byte
// 	question              []byte
// 	options               [][]byte
// 	duration              int
// 	participants          [][]byte
// 	pollID                []byte
// 	fee                   int
// 	signature             []byte
// 	homomorphicPublicKey  *paillier.PublicKey
// 	homomorphicPrivateKey *paillier.PrivateKey
// }

// // CreatePoll creates a new poll with the given parameters.
// func CreatePoll(creatorPublicKey []byte, creatorPrivateKey []byte, question []byte, options [][]byte, duration int, participants [][]byte) *Poll {
// 	p := &Poll{}
// 	p.creatorPublicKey = creatorPublicKey
// 	p.question = question
// 	p.options = options
// 	p.duration = duration
// 	p.participants = participants

// 	// calculate poll ID
// 	p.calculatePollID()

// 	// generate homomorphic key pair
// 	p.generateHomomorphicKeyPair()

// 	// calculate fee
// 	p.calculateFee()

// 	// calculate signature
// 	p.signature = ed25519.Sign(creatorPrivateKey, p.Hash())

// 	return p
// }

// // calculateFee calculates the fee for the poll.
// func (p *Poll) calculateFee() {
// 	questionSize := len(p.question)
// 	participantsSize := len(p.participants)
// 	p.fee = questionSize + participantsSize*256 + p.duration
// }

// // calculatePollID calculates the ID for the poll.
// func (p *Poll) calculatePollID() {
// 	h := sha256.New()
// 	h.Write([]byte(p.question))
// 	var optionStrings []string
// 	for _, option := range p.options {
// 		optionStrings = append(optionStrings, string(option))
// 	}
// 	h.Write([]byte(strings.Join(optionStrings, "")))
// 	var participantStrings []string
// 	for _, participant := range p.participants {
// 		participantStrings = append(participantStrings, string(participant))
// 	}
// 	h.Write([]byte(strings.Join(participantStrings, "")))
// 	p.pollID = h.Sum(nil)
// }

// // generateHomomorphicKeyPair generates the homomorphic key pair for the poll.
// func (p *Poll) generateHomomorphicKeyPair() {
// 	privKey, err := paillier.GenerateKey(rrand.Reader, 2048)
// 	if err != nil {
// 		panic(err)
// 	}
// 	p.homomorphicPrivateKey = privKey
// 	p.homomorphicPublicKey = &privKey.PublicKey
// }

// // Hash calculates the hash of the poll.
// func (p *Poll) Hash() []byte {
// 	h := sha256.New()
// 	h.Write(p.GetPoll())
// 	return h.Sum(nil)
// }

// // GetPoll returns the poll data as a byte array.
// func (p *Poll) GetPoll() []byte {
// 	data := map[string]interface{}{
// 		"creatorPublicKey":     p.creatorPublicKey,
// 		"question":             p.question,
// 		"homomorphicPublicKey": p.homomorphicPublicKey,
// 		"options":              p.options,
// 		"duration":             p.duration,
// 		"participants":         p.participants,
// 		"pollID":               p.pollID,
// 		"fee":                  p.fee,
// 		"signature":            p.signature,
// 	}
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return jsonData
// }

// // LoadPoll deserializes a poll from JSON data.
// func LoadPoll(data []byte) (*Poll, error) {
// 	p := &Poll{}
// 	var poll map[string]interface{}
// 	err := json.Unmarshal(data, &poll)
// 	if err != nil {
// 		return nil, err
// 	}
// 	p.creatorPublicKey = poll["creatorPublicKey"].([]byte)
// 	p.homomorphicPublicKey = poll["homomorphicPublicKey"].(*paillier.PublicKey)
// 	p.question = poll["question"].([]byte)
// 	p.options = poll["options"].([][]byte)
// 	p.duration = int(poll["duration"].(float64))
// 	p.participants = poll["participants"].([][]byte)
// 	p.pollID = poll["pollID"].([]byte)
// 	p.fee = int(poll["fee"].(float64))
// 	p.signature = poll["signature"].([]byte)

// 	return p, nil
// }

// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// func RandStringBytes(n int) []byte {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}
// 	return b
// }

// func findMoney() int {
// 	return 10000000000000
// }
// func calculateHash(p Poll) []byte {
// 	h := sha256.New()
// 	h.Write(p.GetPoll())
// 	return h.Sum(nil)
// }

// func verifyPoll(poll *Poll) bool {
// 	money := findMoney()
// 	if poll.fee > money {
// 		fmt.Println("Not enough blocked money for the user")
// 		return false
// 	}
// 	expectedHash := calculateHash(*poll)
// 	if string(poll.Hash()) != string(expectedHash) {
// 		fmt.Println("Poll hash mismatch")
// 		return false
// 	}
// 	return true
// }

// type NIZKProof struct {
// 	c, z *big.Int
// }

// type Vote struct {
// 	poll *Poll

// 	vote          []byte
// 	encryptedVote [][]byte

// 	wallet     *wallet.Wallet
// 	privateKey *ecdsa.PrivateKey

// 	publicWitnessBuff []byte
// 	vkBuf             []byte
// 	proofBuf          []byte

// 	signature []byte

// 	rand *big.Int
// }

// // - set poll ID
// // - set vote
// // - encrypt vote with public key of poll creator
// // - create a one-time key pair
// // - create proof of authenticity (that one-time key pair is generated by master keys)
// // - create proof of authorization (that master public key is in participants list)
// // - calculate key image
// // - sign vote with private spend key
// func (v *Vote) CreateVote(wallet wallet.Wallet, pollData []byte, vote []byte) {
// 	v.wallet = &wallet
// 	poll, err := LoadPoll(pollData)
// 	if err != nil {
// 		panic(err)
// 	}
// 	v.poll = poll
// 	v.vote = vote
// 	v.rand = randomBigInt(v.poll.homomorphicPublicKey.N)

// 	v.calculateEncryptedVote()

// 	v.privateKey = wallet.GenerateOneTimePair(v.rand.Bytes())

// 	v.calculateProof()

// }

// func randomBigInt(max *big.Int) *big.Int {
// 	r, err := rrand.Int(rrand.Reader, max)
// 	if err != nil {
// 		log.Fatalf("Failed to generate random number: %v", err)
// 	}
// 	return r
// }

// func (v *Vote) calculateEncryptedVote() {
// 	for _, dataVote := range v.vote {
// 		enc, err := paillier.EncryptWithNonce(v.poll.homomorphicPublicKey, v.rand, []byte{dataVote})
// 		if err != nil {
// 			panic(err)
// 		}
// 		v.encryptedVote = append(v.encryptedVote, enc.Bytes())
// 	}
// }

// type Circuit struct {
// 	Random        frontend.Variable `gnark:",public"`
// 	PublicKey     eddsa.PublicKey
// 	ListPublicKey []eddsa.PublicKey `gnark:",public"`
// 	Signature     eddsa.Signature   `gnark:",public"`
// }

// func (circuit *Circuit) Define(api frontend.API) error {

// 	curve, err := twistededwards.NewEdCurve(api, 1)
// 	if err != nil {
// 		return err
// 	}

// 	hash, err := mimc.NewMiMC(api)
// 	if err != nil {
// 		return err
// 	}

// 	eddsa.Verify(curve, circuit.Signature, circuit.Random, circuit.PublicKey, &hash)

// 	temp := frontend.Variable(0)
// 	for i := 0; i < len(circuit.ListPublicKey); i++ {
// 		temp = api.Add(temp, api.Select(
// 			api.And(
// 				api.IsZero(api.Cmp(circuit.ListPublicKey[i].A.X, circuit.PublicKey.A.X)),
// 				api.IsZero(api.Cmp(circuit.ListPublicKey[i].A.Y, circuit.PublicKey.A.Y)),
// 			), frontend.Variable(1), frontend.Variable(0)))
// 	}
// 	api.AssertIsEqual(temp, frontend.Variable(1))
// 	return nil
// }

// func (v *Vote) calculateProof() {
// 	circuit := Circuit{
// 		ListPublicKey: make([]eddsa.PublicKey, 3),
// 	}

// 	// Compile the circuit
// 	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
// 	if err != nil {
// 		fmt.Println("Circuit compilation error:", err)
// 		return
// 	}

// 	fmt.Println("Circuit compiled successfully", ccs)

// 	// Create the Prover and Verifier keys, then generate the proof and verify

// 	// Setup the proving and verifying keys (trusted setup)
// 	pk, vk, err := groth16.Setup(ccs)
// 	if err != nil {
// 		fmt.Println("Setup error:", err)
// 		return
// 	}
// 	fmt.Println("Setup successful")

// 	// witness definition
// 	msg := v.privateKey.PublicKey.Bytes()

// 	signature, err := v.wallet.Sign(msg)
// 	if err != nil {
// 		fmt.Println("Error signing message:", err)
// 		return
// 	}

// 	// declare the witness
// 	var assignment Circuit

// 	// assign message value
// 	assignment.Random = msg

// 	// assign public key values
// 	assignment.PublicKey.Assign(1, v.wallet.GetPublicKey().Bytes()[:32])

// 	// assign signature values
// 	assignment.Signature.Assign(1, signature)

// 	var publicKey1 eddsa.PublicKey
// 	s1 := make([]byte, 32)
// 	rand.Read(s1)
// 	publicKey1.Assign(1, s1)
// 	var publicKey2 eddsa.PublicKey
// 	s2 := make([]byte, 32)
// 	rand.Read(s2)
// 	publicKey2.Assign(1, s2)
// 	var publicKey3 eddsa.PublicKey
// 	s3 := make([]byte, 32)
// 	rand.Read(s3)
// 	publicKey3.Assign(1, s3)

// 	assignment.ListPublicKey = []eddsa.PublicKey{publicKey1, publicKey2, assignment.PublicKey}
// 	fmt.Println(assignment)

// 	witness1, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
// 	publicWitness, _ := witness1.Public()
// 	m, _ := publicWitness.MarshalBinary()
// 	fmt.Println("Witness:", string(m))

// 	// Generate a proof
// 	proof, err := groth16.Prove(ccs, pk, witness1)
// 	if err != nil {
// 		fmt.Println("Proof generation error:", err)
// 		return
// 	}

// 	var proofBuf bytes.Buffer
// 	var proofWriter io.Writer = &proofBuf
// 	proof.WriteTo(proofWriter)

// 	var vkBuf bytes.Buffer
// 	var vkWriter io.Writer = &vkBuf
// 	vk.WriteTo(vkWriter)

// 	publicWitnessBuff, err := publicWitness.MarshalBinary()
// 	if err != nil {
// 		fmt.Println("Marshalling error:", err)
// 		return
// 	}

// 	v.publicWitnessBuff = publicWitnessBuff
// 	v.vkBuf = vkBuf.Bytes()
// 	v.proofBuf = proofBuf.Bytes()

// }

// func (v *Vote) GetHash() []byte {
// 	h := sha256.New()
// 	h.Write(v.poll.pollID)
// 	for _, encVote := range v.encryptedVote {
// 		h.Write(encVote)
// 	}
// 	h.Write(v.privateKey.PublicKey.Bytes())
// 	h.Write(v.publicWitnessBuff)
// 	h.Write(v.signature)
// 	h.Write(v.rand.Bytes())
// 	return h.Sum(nil)
// }

// func (v *Vote) GetVote() []byte {
// 	return v.vote
// }

// func LoadVote(data []byte) *Vote {
// 	// vote := new(Vote)
// 	// vote.poll = new(Poll)
// 	// vote.poll.LoadPoll(data)
// 	return nil
// }

// func run(length, sizeOfParticipants int) {
// 	// Measure wall clock time and memory before function execution
// 	// Generate a new key pair using crypto/ed25519
// 	pubKey, privKey, err := ced.GenerateKey(nil)
// 	if err != nil {
// 		fmt.Println("Error generating key pair:", err)
// 		return
// 	}

// 	allParticipants := 0
// 	minSizeOfQuestions := 50
// 	maxSizeOfQuestions := 150
// 	allOptions := 0
// 	minSizeOfOptions := 2
// 	maxSizeOfOptions := 5
// 	minOptions := 10
// 	maxOptions := 50
// 	maxDuration := 432
// 	listOfPolls := make([]Poll, length)

// 	listOfQuestions := make([][]byte, length)
// 	for i := 0; i < length; i++ {
// 		listOfQuestions[i] = RandStringBytes(rand.Intn(maxSizeOfQuestions-minSizeOfQuestions) + minSizeOfQuestions)
// 	}

// 	listOfOptions := make([][][]byte, length)
// 	for i := 0; i < length; i++ {
// 		sizeOfOptions := rand.Intn(maxSizeOfOptions-minSizeOfOptions) + minSizeOfOptions
// 		listOfOptions[i] = make([][]byte, sizeOfOptions)
// 		allOptions += sizeOfOptions
// 		for j := 0; j < sizeOfOptions; j++ {
// 			listOfOptions[i][j] = RandStringBytes(rand.Intn(maxOptions-minOptions) + minOptions)
// 		}
// 	}

// 	listOfParticipants := make([][][]byte, length)
// 	for i := 0; i < length; i++ {

// 		listOfParticipants[i] = make([][]byte, sizeOfParticipants)
// 		allParticipants += sizeOfParticipants
// 		for j := 0; j < sizeOfParticipants; j++ {
// 			listOfParticipants[i][j] = RandStringBytes(64)
// 		}
// 	}

// 	start := time.Now()
// 	var memStart, memEnd runtime.MemStats
// 	runtime.ReadMemStats(&memStart)

// 	// Create a list of polls
// 	for i := 0; i < length; i++ {
// 		listOfPolls[i] = *CreatePoll(pubKey, privKey, listOfQuestions[i], listOfOptions[i], rand.Intn(maxDuration), listOfParticipants[i])
// 	}

// 	// Measure wall clock time and memory after function execution
// 	elapsed := time.Since(start)
// 	runtime.ReadMemStats(&memEnd)

// 	// Calculate the memory usage
// 	memUsage := memEnd.Alloc - memStart.Alloc

// 	size := 0
// 	for i := 0; i < length; i++ {
// 		size += len(listOfPolls[i].GetPoll())
// 	}

// 	fmt.Printf("Poll creation Execution Time: %s\n", elapsed)
// 	fmt.Printf("Poll creation Memory Usage: %d bytes\n", memUsage)
// 	fmt.Printf("size of all Options: %d \n", allOptions)
// 	fmt.Printf("Size of all Participants: %d \n", allParticipants)
// 	fmt.Printf("Size of Polls: %d \n", len(listOfPolls))
// 	fmt.Printf("Size of poll: %d  byte\n", size)
// 	fmt.Println("----------------------")

// 	startMine := time.Now()
// 	var memStartMine, memEndMine runtime.MemStats
// 	runtime.ReadMemStats(&memStartMine)

// 	// Create a list of polls
// 	for i := 0; i < length; i++ {
// 		verifyPoll(&listOfPolls[i])
// 	}

// 	// Measure wall clock time and memory after function execution
// 	elapsedMine := time.Since(startMine)
// 	runtime.ReadMemStats(&memEndMine)

// 	// Calculate the memory usage
// 	memUsageMine := memEndMine.Alloc - memStartMine.Alloc

// 	fmt.Printf("mining Execution Time: %s\n", elapsedMine)
// 	fmt.Printf("mining Memory Usage: %d bytes\n", memUsageMine)

// 	fmt.Println("********************************")
// }

// func main() {
// 	times := []int{1, 2, 3, 4, 5}
// 	for _, j := range times {
// 		fmt.Println("||||||||||||||||||||||||||||")
// 		fmt.Println("||||||||||||||||||||||||||||")
// 		fmt.Println("||||||||||||||||||||||||||||")
// 		fmt.Println("||||||||||||||||||||||||||||")
// 		fmt.Println(j)

// 		listparticipants := []int{1, 2, 5, 10, 15, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 100000, 1000000}
// 		for _, i := range listparticipants {
// 			fmt.Println("|||||||||")
// 			fmt.Println(i)
// 			run(1, i)
// 			run(5, i)
// 			run(10, i)
// 			run(25, i)
// 			run(50, i)
// 			run(100, i)
// 			run(250, i)
// 			run(500, i)
// 			run(1000, i)
// 			run(5000, i)
// 		}
// 	}
// }
