package main

import (
	"encoding/json"
	"fmt"
	"github.com/pburskey/hasd_covid/domain"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

func main() {

	school := "HES"
	category := "Students"
	var schoolMetrics []string = getMetricsForSchoolAndCategory(school, category)

	details := make([]domain.DataPoint, 0)

	for _, metric := range schoolMetrics {
		metricDetailJson := getMetricsDetail(metric)
		details = append(details, metricDetailJson)
	}

	sort.Slice(details, func(i int, j int) bool {
		return details[i].DateTime.Before(details[j].DateTime)
	})

	for _, metric := range details {
		log.Println(fmt.Sprintf("\t\tDate:%s\t\tActiveCases: %d\tTotalPositiveCases: %d\tProbableCases: %d\tResolvedCases: %d", metric.DateTime, metric.Metric.ActiveCases, metric.Metric.TotalPositiveCases, metric.Metric.ProbableCases, metric.Metric.ResolvedCases))

	}

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

func getMetricsDetail(aMetricKey string) domain.DataPoint {

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

	var responseObject domain.DataPoint
	json.Unmarshal(responseData, &responseObject)

	return responseObject

}
