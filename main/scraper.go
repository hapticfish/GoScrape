package main

import (
	"fmt"
	"io"
	"net/http"
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

	// Process the HTML
	doc.Find(".main-news-controls__wrap > a").Each(func(i int, s *goquery.Selection) {

		text := s.Text()

		mutex.Lock()
		_, err := writer.Write([]byte(text + "\n"))
		mutex.Unlock()

		if err != nil {
			fmt.Println(err)
		}

	})
}
