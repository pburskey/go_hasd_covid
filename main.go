package main

import (
	"fmt"
	"log"
	"net/http"
	"golang.org/x/net/html"

)

func main() {

	covidAsString, err := Extract()
	if err != nil{
		log.Fatal(fmt.Errorf("Unable to acquire covid data: %s: %v", covidAsString, err))
	}

}

var hasd_URL string  = "https://www.hasd.org/community/covid-19-daily-updates.cfm"

func Extract() (covidAsString string, err error){
	resp , err := http.Get(hasd_URL)
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



	covidAsString = doc.Data

	return
}
