package scraper

import (
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type ScrapedItem struct {
	Title        string
	Price        string
	Availability string
}

func ScrapeWebsite(baseURL string) ([]ScrapedItem, error) {

	c := colly.NewCollector(
		//Spread out allowed urls entries as parameters
		//AllowedDomains is a domain whitelist
		colly.AllowedDomains("books.toscrape.com"),

		// Set User-Agent to mimic a browser, otherwise I receive forbidden status code
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),

		//Enabling on Asynchronous Requests (need to set limits after this outwith the collector instance, also need to set Wait() to ensure all requests are finished)
		// colly.Async(true),
	)
	//Setting limits for the collector
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	//OnRequest – runs when the program sends a request to the server.
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL, "Response pending...")
	})

	//OnError – runs when or if we receive an error from the server. In Colly, this is any response that isn’t in the 200’s for server codes.
	c.OnError(func(r *colly.Response, err error) {
		if visitErr := c.Visit(baseURL); visitErr != nil {
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
		//Visit the linked page to scrape books within the category
		e.Request.Visit(e.Request.AbsoluteURL(categoryLink))
	})

	//OnHTML – runs when the program accesses the HTML resource that was served to it.
	//Looking for product_pod class (this has child elements that include the data we need for the first page)

	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		title := e.ChildAttr("div img", "alt")
		price := e.ChildText("p.price_color")
		availability := e.ChildText("p.instock.availability")
		//Get the category from the URL (assuming it follows the pattern "category/{category_name}/index.html")
		fmt.Printf("Title: %s\nPrice: %s\nAvailability: %s\n", title, price, availability)
		//Add the data to the slice
		item := ScrapedItem{Title: title, Price: price, Availability: availability}
		scrapedData = append(scrapedData, item)
	})

	log.Println("Starting crawl at: ", baseURL)

	//Wait for all requests to finish
	c.Wait()

	if err := c.Visit(baseURL); err != nil {
		return nil, err
	}

	//Return the scraped data
	return scrapedData, nil
}
