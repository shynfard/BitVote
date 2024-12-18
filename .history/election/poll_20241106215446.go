package election

import (
	"fmt"
	"math/rand"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/ecdsa"

	rrand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"strings"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"github.com/shynfard/BitVote/blockchain"
	"github.com/shynfard/BitVote/wallet"
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

// Poll represents an election poll.
type PollPublic struct {
	creatorPublicKey []byte   `json:"creatorPublicKey"`
	question         []byte   `json:"question"`
	options          [][]byte `json:"options"`
	duration         int      `json:"duration"`
	participants     [][]byte `json:"participants"`
	pollID           []byte   `json:"pollID"`
	fee              int      `json:"fee"`
	signature        []byte   `json:"signature"`
}

// CreatePoll creates a new poll with the given parameters.
func CreatePoll(creatorPrivateKey []byte, question []byte, options [][]byte, duration int, participants [][]byte) *Poll {
	p := &Poll{}
	priv := ecdsa.PrivateKey{}
	priv.SetBytes(creatorPrivateKey)
	p.creatorPublicKey = priv.Public().Bytes()
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
	hash := p.Hash()
	signature, err := priv.Sign(hash, sha256.New())
	if err != nil {
		panic(err)
	}
	p.signature = signature

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
	h.Write(p.GetPollForHash())
	return h.Sum(nil)
}

// GetPoll returns the poll data as a byte array.
func (p *Poll) GetPoll() []byte {
	// pollPublic := map[string]interface{}{
	// 	"creatorPublicKey":     p.creatorPublicKey,
	// 	"question":             p.question,
	// 	"homomorphicPublicKey": p.homomorphicPublicKey,
	// 	"options":              p.options,
	// 	"duration":             p.duration,
	// 	"participants":         p.participants,
	// 	"pollID":               p.pollID,
	// 	"fee":                  p.fee,
	// 	"signature":            p.signature,
	// }
	pollPublic := PollPublic{
		creatorPublicKey: p.creatorPublicKey,
		question:         p.question,
		options:          p.options,
		duration:         p.duration,
		participants:     p.participants,
		pollID:           p.pollID,
		fee:              p.fee,
		signature:        p.signature,
	}
	fmt.Println(pollPublic)
	jsonData, err := json.Marshal(pollPublic)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// GetPoll returns the poll data as a byte array.
func (p *Poll) GetPollForHash() []byte {

	data := map[string]interface{}{
		"creatorPublicKey":     p.creatorPublicKey,
		"question":             p.question,
		"homomorphicPublicKey": p.homomorphicPublicKey,
		"options":              p.options,
		"duration":             p.duration,
		"participants":         p.participants,
		"pollID":               p.pollID,
		"fee":                  p.fee,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// LoadPoll deserializes a poll from JSON data.
func LoadPoll(data []byte) (*PollPublic, error) {
	var poll PollPublic
	err := json.Unmarshal(data, &poll)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func VerifyPoll(poll *Poll) bool {
	// Check money
	money := blockchain.GetMoneyByPublicKey(poll.creatorPublicKey)
	if poll.fee > money {
		fmt.Println("Not enough blocked money for the user")
		return false
	}
	// verify signature
	pubKey := ecdsa.PublicKey{}
	pubKey.SetBytes(poll.creatorPublicKey)
	ok, err := pubKey.Verify(poll.signature, poll.Hash(), sha256.New())
	if !ok || err != nil {
		fmt.Println(poll.Hash())
		return false
	}
	return true
}
