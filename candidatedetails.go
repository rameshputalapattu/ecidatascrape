package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type CandidateDetail struct {
	bun.BaseModel    `bun:"table:candidate_detail"`
	Name             string
	Party            string
	NumVotes         int64
	ConstituencyCode string
}

func GetAllCandidateDetails(constituencies []Constituency) []CandidateDetail {
	var candidatedetails []CandidateDetail
	for _, constituency := range constituencies {
		candidatedetails = append(candidatedetails, GetCandidateDetails(constituency.Code)...)
	}
	return candidatedetails

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
			if numVotesStr != "Uncontested" {
				log.Fatal("error converting str to int:", err, constituencyCode)
			}

		}

		candName, err := item.Locator(".nme-prty > h5").TextContent()
		if err != nil {
			log.Fatal("unable to fetch candidate name:", err)
		}
		party, err := item.Locator(".nme-prty > h6").TextContent()
		if err != nil {
			log.Fatal("unable to fetch party name:", party)
		}

		candidatedetail := CandidateDetail{
			Name:             candName,
			Party:            party,
			NumVotes:         int64(numVotes),
			ConstituencyCode: constituencyCode,
		}
		candidatedetails = append(candidatedetails, candidatedetail)
	}

	return candidatedetails
}

func loadCandidateDetails(dbFile string, candidatedetails []CandidateDetail) error {

	ctx := context.Background()

	sqldb, err := sql.Open(sqliteshim.ShimName, dbFile)
	if err != nil {
		return err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	drop_table_query := `drop table if exists candidate_detail;`
	_, err = db.ExecContext(ctx, drop_table_query)
	if err != nil {
		return err
	}
	var create_table_query string = `create table if not exists candidate_detail
	(name text,party text,num_votes integer,constituency_code text)`
	_, err = db.ExecContext(ctx, create_table_query)
	if err != nil {
		return err

	}

	_, err = db.NewInsert().Model(&candidatedetails).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
