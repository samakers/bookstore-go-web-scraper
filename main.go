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

	data, err := scraper.ScrapeWebsite("https://books.toscrape.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	//Connect to db and store data
	// db.ConnectToDB()
	//Define and ensure I am passing in correct type here
	dbManager, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Assuming items is your slice of scraped items
	dbManager.StoreInDB(data)

	//Pass data to secondary data store (typesense)

	// client := typesense.CreateTypesenseClient()
	// typesense.IndexInTypesense(client, data)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Total time taken: %s\n", duration)
	fmt.Println("Data:", data)

}
