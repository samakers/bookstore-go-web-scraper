package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type ScrapedItem struct {
	Title        string
	Price        string
	Availability string
}

func main() {

	// Tracking execution time
	startTime := time.Now()
	//Defining base url to scrape
	baseURL := "books.toscrape.com"
	//Concatenate protocol to base URL
	startingURL := "https://" + baseURL
	//Init slice of strings (could allow more than one URL, easy to change)
	allowedUrls := []string{baseURL}
	//Colly’s main entity is a Collector object.
	//Collector manages the network communication and is responsible for the execution of the attached callbacks while a collector job is running.
	//To work with colly, you have to initialize a Collector:
	c := colly.NewCollector(
		//Spread out allowed urls entries as parameters
		//AllowedDomains is a domain whitelist
		colly.AllowedDomains(allowedUrls...),
		//Enabling on Asynchronous Requests (need to set limits after this outwith the collector instance, also need to set Wait() to ensure all requests are finished)
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	//OnRequest – runs when the program sends a request to the server.
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL, "Response pending...")
	})

	//OnError – runs when or if we receive an error from the server. In Colly, this is any response that isn’t in the 200’s for server codes.
	c.OnError(func(r *colly.Response, err error) {
		if visitErr := c.Visit(startingURL); visitErr != nil {
			log.Println("Something went wrong:", visitErr)
			log.Println("Broken link:", r.Request.URL)
			os.Exit(1)
		}
	})

	//OnResponse – runs when the program receives a response from the server.
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL, "Status Code:", r.StatusCode)
	})

	var scrapedData []ScrapedItem
	//Selector to go through sidebar links for each category
	c.OnHTML("#default > div > div > div > aside > div.side_categories > ul > li > ul > li > a", func(e *colly.HTMLElement) {
		categoryLink := e.Attr("href")
		// visit the linked page to scrape books within the category
		e.Request.Visit(e.Request.AbsoluteURL(categoryLink))
	})

	// /OnHTML – runs when the program accesses the HTML resource that was served to it.
	// looking for product_pod class (this has child elements that include the data we need for the first page)

	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		title := e.ChildAttr("div img", "alt")
		price := e.ChildText("p.price_color")
		availability := e.ChildText("p.instock.availability")
		// Get the category from the URL (assuming it follows the pattern "category/{category_name}/index.html")
		fmt.Printf("Title: %s\nPrice: %s\nAvailability: %s\n", title, price, availability)
		// Add the data to the slice
		item := ScrapedItem{Title: title, Price: price, Availability: availability}
		scrapedData = append(scrapedData, item)
	})

	log.Println("Starting crawl at: ", startingURL)

	if err := c.Visit(startingURL); err != nil {
		log.Println("Error on start of crawl: ", err)
	}

	//Wait for all requests to finish
	c.Wait()
	//Call function to save JSON
	saveJSON("scraped_data.json", scrapedData)
	//finish executing
	endTime := time.Now()
	//calculate final time
	duration := endTime.Sub(startTime)
	fmt.Printf("Total time taken: %s\n", duration)
}

// saveJSON writes the provided data to a JSON file
func saveJSON(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating JSON file:", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		log.Fatal("Error encoding JSON:", err)
	}
	log.Println("JSON data saved to", filename)
}
