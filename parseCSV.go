package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func parseCSV(fileName string) {

	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file: %s for processing", fileName, err)
	}
	r := csv.NewReader(csvFile)
	if r != nil {

		dataMap := make(map[string]map[string]CovidMetric)
		//categories["student"] = make(map[string]CovidMetric)
		//categories["staff"] = make(map[string]CovidMetric)
		//var metric CovidMetric
		//metric.ActiveCases =1
		//metric.ProbableCases = 2
		//metric.ResolvedCases = 3
		//metric.TotalPositiveCases = 4
		//categories["student"]["hhs"] = metric
		//

		var categories Stack
		var schools Stack
		first := true
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			} else if first {
				for _, aString := range record {
					if len(aString) > 0 {
						schools.Push(aString)
					}

				}
				first = false
			} else {
				parseCSVRecord(record, dataMap, &categories, &schools)
			}

		}
	}

}

/*
	students
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)
	staff
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)




*/

type CovidMetric struct {
	ActiveCases        int
	TotalPositiveCases int
	ProbableCases      int
	ResolvedCases      int
}

func parseCSVRecord(record []string, dataMap map[string]map[string]CovidMetric, categories *Stack, schools *Stack) {
	if record == nil || len(record) <= 0 {
		return
	}

	if countNonEmptyValuesIn(record) == 2 {
		var category string = record[1]
		_, ok := dataMap[category]
		if !ok {
			dataMap[category] = make(map[string]CovidMetric)
			for _, aSchool := range schools {

			}
		}
		fmt.Printf("Found category %s\n", category)
		categories.Push(category)
	} else {
		category, _ := categories.Peak()
		var metricName string = record[1]
		fmt.Printf("Category %s Metric: %s ... Record %s\n", category, metricName, record)
		if metricName == "Active Cases" {

		} else if metricName == "Total Positive Cases" {

		} else if metricName == "Probable Cases" {

		} else if metricName == "Resolved (no longer active)" {

		}

	}

}
