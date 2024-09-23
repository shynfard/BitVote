package election

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"strings"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

type Poll struct {
	creatorPublicKey []byte
	// creatorPrivateKey []byte
	question     []byte
	options      [][]byte
	duration     int
	participants [][]byte
	pollID       []byte
	fee          int
	signature    []byte

	homomorphicPublicKey  *paillier.PublicKey
	homomorphicPrivateKey *paillier.PrivateKey
}

func (p *Poll) CreatePoll(creatorPublicKey []byte, creatorPrivateKey []byte, question []byte, options [][]byte, duration int, participants [][]byte) {
	p.creatorPublicKey = creatorPublicKey
	// p.creatorPrivateKey = creatorPrivateKey
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

// fee calculation
func (p *Poll) calculateFee() {
	p.fee = 0
}

// poll ID calculation
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

// generate homomorphic key pair
func (p *Poll) generateHomomorphicKeyPair() {
	// generate private key
	// generate public key
	privKey, err := paillier.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	p.homomorphicPrivateKey = privKey
	p.homomorphicPublicKey = &privKey.PublicKey
}

func (p *Poll) Hash() []byte {
	h := sha256.New()
	h.Write(p.GetPoll())
	return h.Sum(nil)
}

// create a json with data and convert it to byte array
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

// deserialize poll from json
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
