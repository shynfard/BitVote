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

	// msg := []byte{160, 174, 182, 123, 16, 48, 222, 30, 222, 234, 56, 44, 218, 20, 158, 252, 240, 146, 34, 191, 132, 150, 83, 255, 197, 245, 49, 9, 61, 242, 10, 139, 31, 107, 112, 93, 39, 23, 63, 88, 132, 254, 39, 232, 100, 15, 95, 131}
	// msg := []byte("Hello, World!")
	// signature, err := w.Sign(msg)
	// if err != nil {
	// 	panic(err)
	// }

	// create a poll
	poll := election.CreatePoll(w.GetPrivateKey().Bytes(), []byte("What is your favorite color?"), [][]byte{[]byte("Red"), []byte("Blue"), []byte("Green")}, 10, [][]byte{w.GetPublicKey().Bytes()})

	pollByte := poll.GetPoll()

	// create a vote
	vote := election.CreateVote(*w, pollByte, []byte{0, 1, 0})

	fmt.Println("Vote: ", vote)

}
