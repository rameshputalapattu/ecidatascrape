package main

import "fmt"

func main() {
	states := GetStates()
	fillConstituencies(states)
	fmt.Println(states)
	candidatedetails := GetCandidateDetails("S0125")
	fmt.Println(candidatedetails)
}
