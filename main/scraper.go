package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func Scrape(url string, writer io.Writer, wg *sync.WaitGroup, mutex *sync.Mutex) {
	fmt.Printf("Starting scraping: %s...\n", url)
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

	if strings.Contains(url, "cointelegraph.com") {
		doc.Find(".main-news-controls__wrap > a").Each(func(i int, s *goquery.Selection) {
			text := wrapText(s.Text(), 80)
			textToWrite := fmt.Sprintf("Cointelegraph: %s\n\n", text)

			mutex.Lock()
			_, err := writer.Write([]byte(textToWrite))
			mutex.Unlock()

			if err != nil {
				fmt.Println(err)
			}
		})
	} else if strings.Contains(url, "coindesk.com") {
		semaphore := make(chan struct{}, 10) // Semaphore to limit to 10 goroutines
		var innerWg sync.WaitGroup

		doc.Find("a.card-imagestyles__CardImageWrapper-sc-1kbd3qh-0.WDSwd").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				semaphore <- struct{}{} // Acquire a token
				innerWg.Add(1)
				go func(url string) {
					defer innerWg.Done()
					scrapeCoinDeskArticle("https://www.coindesk.com"+url, writer, mutex)
					<-semaphore // Release the token
				}(href)
			}
		})

		innerWg.Wait()   // Wait for all scraping goroutines to finish
		close(semaphore) // Close the semaphore channel
	}

	fmt.Printf("Finished scraping: %s...\n", url)
}

func scrapeCoinDeskArticle(articleURL string, writer io.Writer, mutex *sync.Mutex) {
	resp, err := http.Get(articleURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var articleBody string
	doc.Find("div.contentstyle__StyledWrapper-sc-g5cdrh-0.gkcZwU.composer-content").Each(func(i int, s *goquery.Selection) {
		articleBody += wrapText(s.Text(), 80) + "\n\n"
	})

	mutex.Lock()
	defer mutex.Unlock()
	_, err = writer.Write([]byte("CoinDesk: URL: " + articleURL + "\n" + indentString(articleBody, "    ") + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func wrapText(text string, lineWidth int) string {
	var wrappedText bytes.Buffer
	currentLineLength := 0

	for _, word := range strings.Fields(text) {
		if currentLineLength+len(word)+1 > lineWidth {
			wrappedText.WriteString("\n")
			currentLineLength = 0
		}
		if currentLineLength > 0 {
			wrappedText.WriteString(" ")
			currentLineLength++
		}
		wrappedText.WriteString(word)
		currentLineLength += len(word)
	}

	return wrappedText.String()
}

func indentString(str, indent string) string {
	var indentedText string
	for _, line := range strings.Split(str, "\n") {
		indentedText += indent + line + "\n"
	}
	return indentedText
}
