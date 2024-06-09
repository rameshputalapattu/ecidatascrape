package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type CandidateDetail struct {
	Name             string
	Party            string
	NumVotes         int64
	ConstituencyCode string
}

func GetCandidateDetails(constituencyCode string) []CandidateDetail {
	var candidatedetails []CandidateDetail

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()
	const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	bc, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(userAgent),
		BypassCSP: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not create browser context: %v", err)
	}
	defer bc.Close()
	page, err := bc.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	resp, err := page.Goto("https://results.eci.gov.in/PcResultGenJune2024/candidateswise-" + constituencyCode + ".htm")

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	_ = resp
	candidateItems, err := page.Locator(".cand-info").All()
	if err != nil {
		log.Fatalf("could not get candidate items: %v", err)
	}
	for _, item := range candidateItems {

		numVotesStr, err := item.Locator("div.status > div:nth-child(2)").TextContent()
		if err != nil {
			log.Fatal("unable to fetch num votes:", err)
		}

		numVotesStr = strings.Split(numVotesStr, " ")[0]
		numVotes, err := strconv.Atoi(numVotesStr)
		if err != nil {
			log.Fatal("error converting str to int")
		}

		candName, err := item.Locator(".nme-prty > h5").TextContent()
		if err != nil {
			log.Fatal("unable to fetch candidate name:", err)
		}
		party, err := item.Locator(".nme-prty > h6").TextContent()
		if err != nil {
			log.Fatal("unable to fetch party name:", party)
		}
		log.Println(candName, party, numVotes)
		candidatedetail := CandidateDetail{
			Name:     candName,
			Party:    party,
			NumVotes: int64(numVotes),
		}
		candidatedetails = append(candidatedetails, candidatedetail)
	}

	return candidatedetails
}
