package main

import (
	rrand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"time"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"github.com/shynfard/BitVote/election"
	"github.com/shynfard/BitVote/wallet"
)

func main() {
	fmt.Print("participants, size of poll")
	poll(1, 1)
	poll(1, 2)
	poll(1, 5)
	poll(1, 10)
	poll(1, 100)
	poll(1, 1000)
	poll(1, 10000)
	poll(1, 100000)
	poll(1, 1000000)
	poll(1, 8000000)
	poll(1, 10000000)
	poll(1, 30000000)

	//
	// for _, j := range times {
	// 	fmt.Println("||||||||||||||||||||||||||||")
	// 	fmt.Println(j)
	// 	listparticipants := []int{1, 10, 100, 1000, 10000}
	// 	for _, i := range listparticipants {
	// 		listPolls := []int{1, 5, 10, 100, 1000, 10000}
	// 		for _, k := range listPolls {
	// 			fmt.Printf("iteration: %d  \n participants %d \n polls %d \n", j, i, k)
	// 			poll(k, i)
	// 		}
	// 	}
	// }

	// times := []int{1}
	// listparticipants := []int{100, 1000, 10000}
	// listPolls := []int{1, 5, 10, 100, 1000, 10000}
	// fmt.Printf("iteration,participants,polls,Poll Execution Time,Poll Memory Usage,Poll size of all Options,Poll Size of all Participants,Poll Size of poll,Mining Execution Time,Mining Memory Usage\n")
	// for _, i := range listparticipants {
	// 	for _, k := range listPolls {
	// 		for _, j := range times {
	// 			fmt.Printf(" %d, %d, %d,", j, i, k)
	// 			poll(k, i)
	// 		}
	// 	}
	// }
	// times := []int{1}
	// listparticipants := []int{3000, 5000, 8000, 100000}
	// listPolls := []int{1}
	// fmt.Printf("iteration,participants,polls,Poll Execution Time,Poll Memory Usage,Poll size of all Options,Poll Size of all Participants,Poll Size of poll,Mining Execution Time,Mining Memory Usage\n")
	// for _, i := range listparticipants {
	// 	for _, k := range listPolls {
	// 		for _, j := range times {
	// 			fmt.Printf(" %d, %d, %d,\n", j, i, k)
	// 			vote(k, i)
	// 		}
	// 	}
	// }
	// cc(1)
	// cc(2)
	// cc(5)
	// cc(10)
	// cc(50)
	// cc(80)
	// cc(200)
	// cc(100)
	// cc(1000)
	// cc(10000)

}

func cc(iteration int) {

	fmt.Printf("%d, ", iteration)

	privKey, err := paillier.GenerateKey(rrand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	votes := make([]byte, iteration)
	for i := 0; i < iteration; i++ {
		votes[i] = byte(rand.Intn(2))
	}
	ev := votesEnc(&privKey.PublicKey, votes)

	start := time.Now()
	CountVotes(ev, &privKey.PublicKey)
	elapsed := time.Since(start)
	fmt.Printf("%s", elapsed)
	fmt.Println()

}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func poll(iteration, sizeOfParticipants int) {
	w := &wallet.Wallet{}

	w.Load("stalder lincoln sou sauers telfore curcio bradway orelu theall emmanuel aubarta oman parsaye capon matias zacharie viviene abeu saidel honebein")

	// ENVS
	allParticipants := 0
	minSizeOfQuestions := 50
	maxSizeOfQuestions := 150
	allOptions := 0
	minSizeOfOptions := 2
	maxSizeOfOptions := 5
	minOptions := 10
	maxOptions := 50
	maxDuration := 432

	listOfPolls := make([]election.Poll, iteration)

	listOfQuestions := make([][]byte, iteration)
	for i := 0; i < iteration; i++ {
		listOfQuestions[i] = RandStringBytes(rand.Intn(maxSizeOfQuestions-minSizeOfQuestions) + minSizeOfQuestions)
	}

	listOfOptions := make([][][]byte, iteration)
	for i := 0; i < iteration; i++ {
		sizeOfOptions := rand.Intn(maxSizeOfOptions-minSizeOfOptions) + minSizeOfOptions
		listOfOptions[i] = make([][]byte, sizeOfOptions)
		allOptions += sizeOfOptions
		for j := 0; j < sizeOfOptions; j++ {
			listOfOptions[i][j] = RandStringBytes(rand.Intn(maxOptions-minOptions) + minOptions)
		}
	}

	listOfParticipants := make([][][]byte, iteration)
	for i := 0; i < iteration; i++ {

		listOfParticipants[i] = make([][]byte, sizeOfParticipants)
		allParticipants += sizeOfParticipants
		for j := 0; j < sizeOfParticipants; j++ {
			listOfParticipants[i][j] = RandStringBytes(64)
		}
	}

	// start := time.Now()
	// var memStart, memEnd runtime.MemStats
	// runtime.ReadMemStats(&memStart)

	for i := 0; i < iteration; i++ {
		listOfPolls[i] = *election.CreatePoll(w.GetPrivateKey().Bytes(), listOfQuestions[i], listOfOptions[i], rand.Intn(maxDuration), listOfParticipants[i])
	}

	// elapsed := time.Since(start)
	// runtime.ReadMemStats(&memEnd)

	// Calculate the memory usage
	// memUsage := memEnd.Alloc - memStart.Alloc

	size := 0
	for i := 0; i < iteration; i++ {
		size += len(listOfPolls[i].GetPoll())
	}

	// fmt.Printf("Poll creation Execution Time: %s\n", elapsed)
	// fmt.Printf("Poll creation Memory Usage: %d bytes\n", memUsage)
	// fmt.Printf("size of all Options: %d \n", allOptions)
	// fmt.Printf("Size of all Participants: %d \n", allParticipants)
	// fmt.Printf("Size of poll: %d  byte\n", size)
	// fmt.Println("----------------------")

	// fmt.Printf("Poll creation Execution Time: %s\n", elapsed)
	// fmt.Printf("Poll creation Memory Usage: %d bytes\n", memUsage)
	// fmt.Printf("size of all Options: %d \n", allOptions)
	// fmt.Printf("Size of all Participants: %d \n", allParticipants)
	// fmt.Printf("Size of poll: %d  byte\n", size)
	// fmt.Println("----------------------")

	// fmt.Printf("Poll creation Execution Time,Poll creation Memory Usage,size of all Options,Size of all Participants,Size of poll\n")
	fmt.Printf("%d,%d\n", allParticipants, size)

	// startMine := time.Now()
	// var memStartMine, memEndMine runtime.MemStats
	// runtime.ReadMemStats(&memStartMine)

	// // Create a list of polls
	// for i := 0; i < iteration; i++ {
	// 	election.VerifyPoll(&listOfPolls[i])
	// }

	// // Measure wall clock time and memory after function execution
	// elapsedMine := time.Since(startMine)
	// runtime.ReadMemStats(&memEndMine)

	// // Calculate the memory usage
	// memUsageMine := memEndMine.Alloc - memStartMine.Alloc

	// // fmt.Printf("mining Execution Time: %s\n", elapsedMine)
	// // fmt.Printf("mining Memory Usage: %d bytes\n", memUsageMine)
	// fmt.Printf("%s,%d\n", elapsedMine, memUsageMine)

}

func vote(iteration, sizeOfParticipants int) {
	w := &wallet.Wallet{}
	wVoter := &wallet.Wallet{}

	w.Load("stalder lincoln sou sauers telfore curcio bradway orelu theall emmanuel aubarta oman parsaye capon matias zacharie viviene abeu saidel honebein")
	wVoter.Load("asghar mamad sou sauers telfore curcio bradway orelu theall emmanuel aubarta oman parsaye capon matias zacharie viviene abeu saidel honebein")

	listOfPolls := make([]election.Poll, iteration)
	participants := make([][]byte, sizeOfParticipants)
	for i := 0; i < sizeOfParticipants; i++ {
		participants[i] = RandStringBytes(32)
	}
	participants[rand.Intn(sizeOfParticipants)] = wVoter.GetPublicKey().Bytes()

	listOfVotes := make([]election.Vote, iteration)

	for i := 0; i < iteration; i++ {
		listOfPolls[i] = *election.CreatePoll(w.GetPrivateKey().Bytes(), []byte("how"), [][]byte{[]byte("good"), []byte("bad")}, 100, participants)
	}

	start := time.Now()
	var memStart, memEnd runtime.MemStats
	runtime.ReadMemStats(&memStart)

	for i := 0; i < iteration; i++ {
		listOfVotes[i] = *election.CreateVote(*wVoter, listOfPolls[i].GetPoll(), []byte{0, 1})
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&memEnd)
	memUsage := memEnd.Alloc - memStart.Alloc
	fmt.Printf("%s,%d,", elapsed, memUsage)

	startMine := time.Now()
	var memStartMine, memEndMine runtime.MemStats
	runtime.ReadMemStats(&memStartMine)

	// Create a list of polls
	for i := 0; i < iteration; i++ {
		listOfVotes[i].VerifyVote()
	}

	// Measure wall clock time and memory after function execution
	elapsedMine := time.Since(startMine)
	runtime.ReadMemStats(&memEndMine)

	// Calculate the memory usage
	memUsageMine := memEndMine.Alloc - memStartMine.Alloc

	// fmt.Printf("mining Execution Time: %s\n", elapsedMine)
	// fmt.Printf("mining Memory Usage: %d bytes\n", memUsageMine)
	fmt.Printf("%s,%d\n", elapsedMine, memUsageMine)
}

func voteEncryption(publicKey *paillier.PublicKey, vote byte) big.Int {
	// Encrypt vote with public key of poll creator
	// encryptedVote := make([]big.Int, len(vote))
	// for i := 0; i < len(vote); i++ {
	// 	r, err := rrand.Int(rrand.Reader, publicKey.N)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	enc, err := paillier.EncryptWithNonce(publicKey, r, []byte{vote[i]})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	encryptedVote[i] = *enc
	// }
	// return encryptedVote
	r, err := rrand.Int(rrand.Reader, publicKey.N)
	if err != nil {
		panic(err)
	}
	enc, err := paillier.EncryptWithNonce(publicKey, r, []byte{vote})
	return *enc
}

func votesEnc(publickey *paillier.PublicKey, votes []byte) []big.Int {

	encResult := make([]big.Int, len(votes))
	for i := 0; i < len(votes); i++ {
		encResult[i] = voteEncryption(publickey, votes[i])
	}
	return encResult
}

// func voteEncryption(publicKey *paillier.PublicKey, votes [][]byte) []big.Int {
// 	// Encrypt vote with public key of poll creator
// 	encryptedVote := make([]big.Int, len(votes))
// 	for i := 0; i < len(votes); i++ {
// 		e []big.Int
// 		for _, dataVote := range votes[i] {
// 			r, err := rrand.Int(rrand.Reader, publicKey.N)
// 			if err != nil {
// 				panic(err)
// 			}
// 			enc, err := paillier.EncryptWithNonce(publicKey, r, []byte{dataVote})
// 			if err != nil {
// 				panic(err)
// 			}
// 			encryptedVote[i] = enc
// 		}
// 	}
// 	return encryptedVote
// }

func CountVotes(encryptedVote []big.Int, homomorphicPublicKey *paillier.PublicKey) big.Int {
	var encryptedSum big.Int
	for i, vote := range encryptedVote {
		if i == 0 {
			encryptedSum = vote
		} else {
			encryptedSum.Mul(&encryptedSum, &vote)
			encryptedSum.Mod(&encryptedSum, homomorphicPublicKey.N)
		}
	}
	return encryptedSum
}
