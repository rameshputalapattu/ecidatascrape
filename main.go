package main

import (
	"log"
)

func main() {
	states := GetStates()
	err := loadStates("data/elections_2024.db", states)
	if err != nil {
		log.Fatal("error loading the states data", err)
	}
	constituencies := getAllConstituencies(states)
	err = loadConstituencies("data/elections_2024.db", constituencies)
	if err != nil {
		log.Fatal("error loading constituencies:", err)
	}
	candidatedetails := GetAllCandidateDetails(constituencies)
	err = loadCandidateDetails("data/elections_2024.db", candidatedetails)
	if err != nil {
		log.Fatal("error loading constituencies:", err)
	}

}
