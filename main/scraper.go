package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery" // GoQuery package for parsing HTML
)

func Scrape(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process the HTML
	doc.Find(".main-news-controls__wrap > a").Each(func(i int, s *goquery.Selection) {
		// Extract data using s.Find or s.Attr etc.
		text := s.Text()
		fmt.Println("Scraped text:", text)
	})
}
