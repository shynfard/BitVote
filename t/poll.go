package t

import (
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	ced "crypto/ed25519"
	rrand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"strings"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"github.com/shynfard/BitVote/old/wallet"
)

// Poll represents an election poll.
type Poll struct {
	Creator               wallet.Wallet
	creatorPublicKey      []byte
	question              []byte
	options               [][]byte
	duration              int
	participants          [][]byte
	pollID                []byte
	fee                   int
	signature             []byte
	homomorphicPublicKey  *paillier.PublicKey
	homomorphicPrivateKey *paillier.PrivateKey
}

// CreatePoll creates a new poll with the given parameters.
func CreatePoll(creatorPublicKey []byte, creatorPrivateKey []byte, question []byte, options [][]byte, duration int, participants [][]byte) *Poll {
	p := &Poll{}
	p.creatorPublicKey = creatorPublicKey
	p.question = question
	p.options = options
	p.duration = duration
	p.participants = participants

	// calculate poll ID
	p.calculatePollID()

	// generate homomorphic key pair
	p.generateHomomorphicKeyPair()

	// calculate fee
	p.calculateFee()

	// calculate signature
	p.signature = ed25519.Sign(creatorPrivateKey, p.Hash())

	return p
}

// calculateFee calculates the fee for the poll.
func (p *Poll) calculateFee() {
	questionSize := len(p.question)
	participantsSize := len(p.participants)
	p.fee = questionSize + participantsSize*256 + p.duration
}

// calculatePollID calculates the ID for the poll.
func (p *Poll) calculatePollID() {
	h := sha256.New()
	h.Write([]byte(p.question))
	var optionStrings []string
	for _, option := range p.options {
		optionStrings = append(optionStrings, string(option))
	}
	h.Write([]byte(strings.Join(optionStrings, "")))
	var participantStrings []string
	for _, participant := range p.participants {
		participantStrings = append(participantStrings, string(participant))
	}
	h.Write([]byte(strings.Join(participantStrings, "")))
	p.pollID = h.Sum(nil)
}

// generateHomomorphicKeyPair generates the homomorphic key pair for the poll.
func (p *Poll) generateHomomorphicKeyPair() {
	privKey, err := paillier.GenerateKey(rrand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	p.homomorphicPrivateKey = privKey
	p.homomorphicPublicKey = &privKey.PublicKey
}

// Hash calculates the hash of the poll.
func (p *Poll) Hash() []byte {
	h := sha256.New()
	h.Write(p.GetPoll())
	return h.Sum(nil)
}

// GetPoll returns the poll data as a byte array.
func (p *Poll) GetPoll() []byte {
	data := map[string]interface{}{
		"creatorPublicKey":     p.creatorPublicKey,
		"question":             p.question,
		"homomorphicPublicKey": p.homomorphicPublicKey,
		"options":              p.options,
		"duration":             p.duration,
		"participants":         p.participants,
		"pollID":               p.pollID,
		"fee":                  p.fee,
		"signature":            p.signature,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return jsonData
}

// LoadPoll deserializes a poll from JSON data.
func LoadPoll(data []byte) (*Poll, error) {
	p := &Poll{}
	var poll map[string]interface{}
	err := json.Unmarshal(data, &poll)
	if err != nil {
		return nil, err
	}
	p.creatorPublicKey = poll["creatorPublicKey"].([]byte)
	p.homomorphicPublicKey = poll["homomorphicPublicKey"].(*paillier.PublicKey)
	p.question = poll["question"].([]byte)
	p.options = poll["options"].([][]byte)
	p.duration = int(poll["duration"].(float64))
	p.participants = poll["participants"].([][]byte)
	p.pollID = poll["pollID"].([]byte)
	p.fee = int(poll["fee"].(float64))
	p.signature = poll["signature"].([]byte)

	return p, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func run(length int) {
	// Measure wall clock time and memory before function execution
	// Generate a new key pair using crypto/ed25519
	pubKey, privKey, err := ced.GenerateKey(nil)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}
	fmt.Printf("Public Key: %x\n", pubKey)
	fmt.Printf("Private Key: %x\n", privKey)

	maxSizeOfParticipants := 1000
	minSizeOfParticipants := 100
	allParticipants := 0
	minSizeOfQuestions := 50
	maxSizeOfQuestions := 150
	allOptions := 0
	minSizeOfOptions := 2
	maxSizeOfOptions := 5
	minOptions := 10
	maxOptions := 50
	maxDuration := 432
	listOfPolls := make([]Poll, length)

	listOfQuestions := make([][]byte, length)
	for i := 0; i < length; i++ {
		listOfQuestions[i] = RandStringBytes(rand.Intn(maxSizeOfQuestions-minSizeOfQuestions) + minSizeOfQuestions)
	}

	listOfOptions := make([][][]byte, length)
	for i := 0; i < length; i++ {
		sizeOfOptions := rand.Intn(maxSizeOfOptions-minSizeOfOptions) + minSizeOfOptions
		listOfOptions[i] = make([][]byte, sizeOfOptions)
		allOptions += sizeOfOptions
		for j := 0; j < sizeOfOptions; j++ {
			listOfOptions[i][j] = RandStringBytes(rand.Intn(maxOptions-minOptions) + minOptions)
		}
	}

	listOfParticipants := make([][][]byte, length)
	for i := 0; i < length; i++ {
		sizeOfParticipants := rand.Intn(maxSizeOfParticipants-minSizeOfParticipants) + minSizeOfParticipants

		listOfParticipants[i] = make([][]byte, sizeOfParticipants)
		allParticipants += sizeOfParticipants
		for j := 0; j < sizeOfParticipants; j++ {
			listOfParticipants[i][j] = RandStringBytes(64)
		}
	}

	start := time.Now()
	var memStart, memEnd runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// Create a list of polls
	for i := 0; i < length; i++ {
		listOfPolls[i] = *CreatePoll(pubKey, privKey, listOfQuestions[i], listOfOptions[i], rand.Intn(maxDuration), listOfParticipants[i])
	}

	// Measure wall clock time and memory after function execution
	elapsed := time.Since(start)
	runtime.ReadMemStats(&memEnd)

	// Calculate the memory usage
	memUsage := memEnd.Alloc - memStart.Alloc

	fmt.Printf("Execution Time: %s\n", elapsed)
	fmt.Printf("Memory Usage: %d bytes\n", memUsage)
	fmt.Printf("Size of allOptions: %d \n", allOptions)
	fmt.Printf("Size of allParticipants: %d \n", allParticipants)
	fmt.Printf("Size of listOfPolls: %d \n", len(listOfPolls))
	fmt.Printf("Size of listOfQuestions: %d \n", len(listOfQuestions))

	fmt.Println("Done")
}

func main() {
	run(1)
	run(10)
	run(100)
	run(1000)
	run(10000)
	run(100000)
}
