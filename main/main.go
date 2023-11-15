package main

import (
	"os"
	"sync"
)

func main() {
	file, err := os.Create("scrapedData.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	urls := []string{"https://cointelegraph.com/", "https://www.coindesk.com/"}

	for _, url := range urls {
		wg.Add(1)
		go Scrape(url, file, &wg, mutex) // Call the exported Scrape function
	}

	wg.Wait() // Wait for all goroutines to finish
}
