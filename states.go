package main

import (
	"log"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type State struct {
	Code           string
	Name           string
	Constituencies []Constituency
}

func GetStates() []State {

	var states []State

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
	resp, err := page.Goto("https://results.eci.gov.in/PcResultGenJune2024/index.htm")

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	_ = resp
	stateOptions := page.Locator("#ctl00_ContentPlaceHolder1_Result1_ddlState > option")
	items, err := stateOptions.All()

	if err != nil {
		log.Fatal("could not get All options:", err)
	}

	log.Println("size of rows:", len(items))
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	for _, item := range items {
		stateCode, err := item.GetAttribute("value")

		if err != nil {
			log.Fatalf("could not get statecode: %v", err)
		}

		if strings.Trim(stateCode, " ") == "" {
			continue
		}
		stateText, err := item.TextContent()
		if err != nil {
			log.Fatalf("could not get statecode: %v", err)
		}
		state := State{
			Code: stateCode,
			Name: stateText,
		}
		states = append(states, state)

	}
	return states

}
