package handlers

import (
	"fmt"
	"net/http"
	"web-crawler/scraper"
)

// HomeHandler handles requests to the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page!")
}

// DataHandler handles requests to the /data endpoint
func DataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := scraper.ScrapeWebsite("https://books.toscrape.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Convert the data to JSON and write it to the response
	fmt.Fprint(w, data)
}
