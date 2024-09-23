package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"strconv"
	"strings"

	"github.com/0xdecaf/zkrp/ccs08"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

func StringToIntArray(A string) []int64 {
	strs := strings.Split(A, " ")
	ary := make([]int64, len(strs))
	for i := range ary {
		v, _ := strconv.Atoi(strs[i])
		ary[i] = int64(v)
	}
	return ary
}

func main() {

	tofind := 6
	set := "1 2 3 4 5"

	argCount := len(os.Args[1:])

	if argCount > 0 {
		tofind, _ = strconv.Atoi(os.Args[1])
	}
	if argCount > 1 {
		set = os.Args[2]
	}
	s := StringToIntArray(set)

	p, _ := ccs08.SetupSet(s)
	fmt.Printf("To find: %d\n", tofind)
	fmt.Printf("Set: %v\n", s)

	r, _ := rand.Int(rand.Reader, bn256.Order)

	proof_out, err := ccs08.ProveSet(int64(tofind), r, p)

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}

	result, err2 := ccs08.VerifySet(&proof_out, &p)

	if err2 != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}

	if result == true {
		fmt.Printf("\nProof that %d is in [%v]\n\n", tofind, set)
	} else {
		fmt.Printf("\nDid not find value in array\n")
	}
	fmt.Printf("Proof: C= %s\n", proof_out.C)
	fmt.Printf("Proof: D= %s\n", proof_out.D)
	fmt.Printf("Proof: V= %s\n", proof_out.V)

}
