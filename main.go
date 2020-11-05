package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
)

func main() {

	covidAsString, err := Extract()
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to acquire covid data: %s: %v", covidAsString, err))
	}

}

var hasd_URL string = "https://www.hasd.org/community/covid-19-daily-updates.cfm"

func ExtractUsingColly() {
	URL := url.URL.Query(hasd_URL).Get("url")
	if URL == "" {
		log.Println("URL is bad")
		return
	}
	log.Println("visiting", URL)

	c := colly.NewCollector(
		colly.AllowedDomains("hasd.org"),
		colly.CacheDir("./cache"),
	)

}

func Extract() (covidAsString string, err error) {
	resp, err := http.Get(hasd_URL)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return covidAsString, fmt.Errorf("getting %s: %s", hasd_URL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return covidAsString, fmt.Errorf("parsing %s as html: %v", hasd_URL, err)
	}
	scrapeDataFrom(doc)

	covidAsString = doc.Data

	return
}

func scrapeDataFrom(aDocument *html.Node) {

}
