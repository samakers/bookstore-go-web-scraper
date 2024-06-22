// main.go
package main

import (
	"fmt"
	"log"
	"time"
	"web-crawler/db"
	"web-crawler/scraper"
)

func main() {

	//Tracking execution time
	startTime := time.Now()

	done := make(chan bool)

	// Scrape website in a separate goroutine
	var data []scraper.ScrapedItem
	var err error
	go func() {
		data, err = scraper.ScrapeWebsite("https://books.toscrape.com")
		if err != nil {
			log.Fatalf("Failed to scrape website: %v", err)
		}

		// Connect to DB and store data
		dbManager, err := db.ConnectToDB()
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}

		dbManager.StoreInDB(data)

		// Signal completion on the channel
		done <- true
	}()

	// Wait for the database operation to complete
	<-done

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Println("Data:", data)
	fmt.Printf("Total time taken: %s\n", duration)
}
