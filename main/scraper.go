package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery" // GoQuery package for parsing HTML
)

func Scrape(url string, writer io.Writer, wg *sync.WaitGroup, mutex *sync.Mutex) {
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

	// Determine the selector based on the URL
	var selector string
	if strings.Contains(url, "cointelegraph.com") {
		selector = ".main-news-controls__wrap > a"
	} else if strings.Contains(url, "coindesk.com") {
		selector = "a.card-imagestyles__CardImageWrapper-sc-1kbd3qh-0.WDSwd"
	}

	// Process the HTML
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		var textToWrite string
		if selector == ".main-news-controls__wrap > a" {
			text := s.Text()
			textToWrite = fmt.Sprintf("Cointelegraph: %s\n", text)
		} else {
			href, exists := s.Attr("href")
			if exists {
				textToWrite = fmt.Sprintf("CoinDesk: URL: %s\n", href)
			}
		}

		if textToWrite != "" {
			mutex.Lock()
			_, err := writer.Write([]byte(textToWrite))
			mutex.Unlock()

			if err != nil {
				fmt.Println(err)
			}
		}

	})
}
