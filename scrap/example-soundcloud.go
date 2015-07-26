package scrap

import (
	"fmt"
	"os"

	"github.com/jloup/html/scraper"
)

func Soundcloud() {
	fmt.Println("starting soundcloud scraper example...")

	s, err := scraper.JSONFileToScraper("testdata/exampledata/soundcloud.json")
	if err != nil {
		fmt.Println("error in scraper.JsonFileToScraper", err)
		return
	}

	// http://www.lagasta.com (music blog)
	f, err := os.Open("testdata/exampledata/soundcloud.html")
	if err != nil {
		fmt.Println("error in os.Open", err)
		return
	}

	items, err := scraper.ScrapHTML(s, f)
	if err != nil {
		fmt.Println("error in scraper.ScrapHTMLn", err)
		return
	}

	for _, item := range items {
		fmt.Printf("Soundcloud type '%s' id %s\n", item["type"], item["id"])
	}
}
