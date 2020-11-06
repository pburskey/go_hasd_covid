package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func main() {

	covidAsString, err := Extract()
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to acquire covid data: %s: %v", covidAsString, err))
	}
	ExtractUsingColly()

}

var hasd_URL string = "https://www.hasd.org/community/covid-19-daily-updates.cfm"

func ExtractUsingColly() {

	c := colly.NewCollector(
		colly.AllowedDomains("hasd.org"),
	)

	c.OnHTML("table[id=\"hasd\"] tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			//talk := Talk{}
			//for each line "tr" do amazing things
			//talks = append(talks, talk)
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	err := c.Visit(hasd_URL)
	if err != nil {
		println(err)
	}

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
