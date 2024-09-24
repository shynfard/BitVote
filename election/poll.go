package election

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"strings"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

// Poll represents an election poll.
type Poll struct {
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
func (p *Poll) CreatePoll(creatorPublicKey []byte, creatorPrivateKey []byte, question []byte, options [][]byte, duration int, participants [][]byte) {
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
}

// calculateFee calculates the fee for the poll.
func (p *Poll) calculateFee() {
	p.fee = 0
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
	privKey, err := paillier.GenerateKey(rand.Reader, 2048)
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
