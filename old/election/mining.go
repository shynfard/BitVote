package election

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func verifyVote(vote Vote) bool {
	expectedHash := calculateHash(vote)
	if vote.GetHash() != expectedHash {
		fmt.Println("Poll hash mismatch")
		return false
	}

	newProof := groth16.NewProof(ecc.BN254)
	newProof.ReadFrom(bytes.NewReader(vote.proofBuf))

	newVk := groth16.NewVerifyingKey(ecc.BN254)
	newVk.ReadFrom(bytes.NewReader(vote.vkBuf))

	newPublicWitness, _ := witness.New(ecc.BN254.ScalarField()) //
	newPublicWitness.UnmarshalBinary(vote.publicWitnessBuff)

	// Verify the proof
	err := groth16.Verify(newProof, newVk, *vote.publicWitness)
	if err != nil {
		fmt.Println("Proof verification failed:", err)
		return false
	} else {
		return true
	}
}

func verifyPoll(poll *Poll) bool {
	money := Blockchain.findMonez(poll.Creator.GetPublicKey())
	if poll.fee < money {
		fmt.Println("Not enough blocked money for the user")
		return false
	}
	expectedHash := calculateHash(poll)
	if poll.Hash != expectedHash {
		fmt.Println("Poll hash mismatch")
		return false
	}
	return true
}

func countVotes(votes []Vote, pubKey *ecdsa.PublicKey) *big.Int {

	encryptedSum := new(big.Int).SetInt64(1)

	for _, vote := range votes {
		if verifyVote(vote) {
			// Add the encrypted vote to the encrypted sum
			encryptedSum.Mul(encryptedSum, vote.encryptedVote)
			encryptedSum.Mod(encryptedSum, pubKey.Params().N)
		}
	}
	return encryptedSum

}
