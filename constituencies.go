package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/playwright-community/playwright-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Constituency struct {
	bun.BaseModel `bun:"table:constituency"`
	Code          string
	Name          string
	StateCode     string
}

func getAllConstituencies(states []State) []Constituency {
	var constituencies []Constituency
	for idx := range states {

		state := states[idx].Code
		constituencies = append(constituencies, getConstituencies(state)...)

	}
	return constituencies

}

func getConstituencies(state string) []Constituency {

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
	resp, err := page.Goto("https://results.eci.gov.in/PcResultGenJune2024/partywiseresult-" + state + ".htm")

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	_ = resp

	constituencyLoc, err := page.Locator("#ctl00_ContentPlaceHolder1_Result1_ddlState > option").All()
	if err != nil {
		log.Fatalf("could not get constituency options: %v", err)
	}

	var constituencies []Constituency
	for _, item := range constituencyLoc {
		constCode, err := item.GetAttribute("value")

		if err != nil {
			log.Fatalf("could not get Constituency code: %v", err)
		}

		if strings.Trim(constCode, " ") == "" {
			continue
		}
		constName, err := item.TextContent()
		if err != nil {
			log.Fatalf("could not get statecode: %v", err)
		}
		constituency := Constituency{
			Code:      constCode,
			Name:      constName,
			StateCode: state,
		}
		constituencies = append(constituencies, constituency)
	}
	//log.Println(constituencies)
	return constituencies

}

func loadConstituencies(dbFile string, constituencies []Constituency) error {

	ctx := context.Background()

	sqldb, err := sql.Open(sqliteshim.ShimName, dbFile)
	if err != nil {
		return err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	drop_table_query := `drop table if exists constituency;`
	_, err = db.ExecContext(ctx, drop_table_query)
	if err != nil {
		return err
	}
	var create_table_query string = `create table if not exists constituency
	(name text,code text,state_code text)`
	_, err = db.ExecContext(ctx, create_table_query)
	if err != nil {
		return err

	}

	_, err = db.NewInsert().Model(&constituencies).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
