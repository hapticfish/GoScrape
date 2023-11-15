package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	urls := []string{"https://cointelegraph.com/", "https://www.coindesk.com/"}

	for _, url := range urls {
		wg.Add(1)
		go Scrape(url, &wg) // Call the exported Scrape function
	}

	wg.Wait() // Wait for all goroutines to finish
}
