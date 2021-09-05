package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pburskey/hasd_covid/domain"
	"github.com/pburskey/hasd_covid/utility"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

func main() {

	arguments := os.Args
	if len(arguments) != 3 {
		log.Fatalln("School and Category arguments are required")
	}

	school := arguments[1]
	category := arguments[2]

	var schoolMetrics []string = getMetricsForSchoolAndCategory(school, category)

	details := make([]domain.RawDataPoint, 0)

	for _, metric := range schoolMetrics {
		metricDetailJson := getMetricsDetail(metric)
		details = append(details, metricDetailJson)
	}

	sort.Slice(details, func(i int, j int) bool {
		return details[i].DateTime.Before(details[j].DateTime)
	})

	var lastMetric domain.RawDataPoint

	reportData := make([][]string, 0)
	for _, currentMetric := range details {

		sign := " "
		signDeterminationFunction := func(a int, b int, skip bool) string {
			sign := " "
			if !skip {
				if a > b {
					sign = "+"
				} else if a < b {
					sign = "-"
				} else {
					sign = " "
				}
			}

			return sign
		}

		row := make([]string, 0)
		row = append(row, utility.AsYYYY_MM_DD_HH24(currentMetric.DateTime))

		skip := lastMetric.Id == 0
		sign = signDeterminationFunction(currentMetric.Metric.ActiveCases, lastMetric.Metric.ActiveCases, skip)
		activeCases := fmt.Sprintf("%s %d", sign, currentMetric.Metric.ActiveCases)
		row = append(row, activeCases)

		sign = signDeterminationFunction(currentMetric.Metric.TotalPositiveCases, lastMetric.Metric.TotalPositiveCases, skip)
		totalPositiveCases := fmt.Sprintf("%s %d", sign, currentMetric.Metric.TotalPositiveCases)
		row = append(row, totalPositiveCases)

		sign = signDeterminationFunction(currentMetric.Metric.ProbableCases, lastMetric.Metric.ProbableCases, skip)
		probableCases := fmt.Sprintf("%s %d", sign, currentMetric.Metric.ProbableCases)
		row = append(row, probableCases)

		sign = signDeterminationFunction(currentMetric.Metric.ResolvedCases, lastMetric.Metric.ResolvedCases, skip)
		resolvedCases := fmt.Sprintf("%s %d", sign, currentMetric.Metric.ResolvedCases)
		row = append(row, resolvedCases)

		reportData = append(reportData, row)
		//log.Println(fmt.Sprintf("\t\tDate:%s\tActiveCases: %s\t\tTotalPositiveCases: %s\t\tProbableCases: %s\t\tResolvedCases: %s", currentMetric.DateTime, activeCases, totalPositiveCases, probableCases, resolvedCases))

		lastMetric = currentMetric
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "ActiveCases", "TotalPositiveCases", "ProbableCases", "ResolvedCases"})
	table.SetBorder(false)       // Set Border to false
	table.AppendBulk(reportData) // Add Bulk Data
	table.Render()

}

func getMetricsForSchoolAndCategory(school string, category string) []string {

	uri := fmt.Sprintf("http://127.0.0.1:8080/api/v1/school/%s/category/%s/metrics", school, category)
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}

	var responseObject []string
	json.Unmarshal(responseData, &responseObject)

	return responseObject

}

func getMetricsDetail(aMetricKey string) domain.RawDataPoint {

	uri := fmt.Sprintf("http://127.0.0.1:8080/api/v1/metric/%s", aMetricKey)
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}

	var responseObject domain.RawDataPoint
	json.Unmarshal(responseData, &responseObject)

	return responseObject

}
