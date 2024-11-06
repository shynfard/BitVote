package main

import (
	"fmt"

	"github.com/shynfard/BitVote/election"
	"github.com/shynfard/BitVote/wallet"
)

func main() {
	// Create a new wallet
	w := &wallet.Wallet{}

	w.Load("stalder lincoln sou sauers telfore curcio bradway orelu theall emmanuel aubarta oman parsaye capon matias zacharie viviene abeu saidel honebein")

	// msg := []byte{160, 174, 182, 123, 16, 48, 222, 30, 222, 234, 56, 44, 218, 20, 158, 174, 182, 123, 16, 48, 222, 30, 222, 234, 56, 44, 218, 20, 158, 174, 182, 123, 16, 48, 222, 30, 222, 234, 56, 44, 218, 20, 158, 234, 56, 44, 218, 20, 158}
	// hFunc := Mim.New()

	// hFunc.Write(msg)
	// hFunc.Write(msg)
	// hFunc.Write(msg)
	// hashed := hFunc.Sum(nil)

	// fmt.Println("Hashed: ", hashed)

	// msg := []byte("Hello, World!")
	// fmt.Println(len(msg))
	// signature, err := w.Sign(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Signature: ", signature)

	// create a poll
	poll := election.CreatePoll(w.GetPrivateKey().Bytes(), []byte("What is your favorite color?"), [][]byte{[]byte("Red"), []byte("Blue"), []byte("Green")}, 10, [][]byte{w.GetPublicKey().Bytes()})

	res := election.VerifyPoll(poll)

	if res != true {
		fmt.Println("Poll verification failed")
	}

	pollByte := poll.GetPoll()

	// create a vote
	vote := election.CreateVote(*w, pollByte, []byte{0, 1, 0})

	fmt.Println("Vote: ", vote)

}
