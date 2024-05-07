// main.go
package main

import (
	"fmt"
	"time"
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
	// database.StoreInDB(data)

	// client := typesense.CreateTypesenseClient()
	// typesense.IndexInTypesense(client, data)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Total time taken: %s\n", duration)
	fmt.Println("Data:", data)

}
