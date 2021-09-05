package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/pburskey/hasd_covid/domain"
	redis_utility "github.com/pburskey/hasd_covid/redis"
	"github.com/pburskey/hasd_covid/utility"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const COVID_DATA = "covid_data_"
const CSV = ".csv"

func parseDateFromFileName(fileName string) time.Time {
	start := strings.Index(fileName, COVID_DATA)

	/*
		strip out the covid_data
	*/
	dateAndTimePortion := fileName[start+len(COVID_DATA) : (len(fileName) - len(CSV))]

	aTime, err := utility.ParseYYYYMMDDHH24MiSS(dateAndTimePortion)
	if err != nil {
		log.Fatal("Unable to parse date and time portion from file name: %s", fileName)
	}
	return aTime
}

func ParseCSV(fileName string, redis *redis_utility.RedisConnection) (time.Time, map[string]map[string]*domain.CovidMetric) {

	key := fmt.Sprintf("file:%s_data", fileName)
	dataMap := make(map[string]map[string]*domain.CovidMetric)
	aTime := parseDateFromFileName(fileName)

	conn := redis.GetRedisConnection()
	defer conn.Close()

	if redis.CacheKeyExists(key) {
		return aTime, nil

	}

	//fmt.Println("Date and time: %s from file name: %s", aTime, fileName)
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file: %s for processing", fileName, err)
	}
	r := csv.NewReader(csvFile)
	if r == nil {
		return aTime, nil
	}

	var categories utility.Stack
	var schools utility.Stack
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
	//log.Println(dataMap)

	aCacheSetResult, err := conn.Do("SET", key, dataMap)
	//log.Print(aCacheSetResult)
	if err != nil || aCacheSetResult != "OK" {
		log.Fatal(err)
	}

	return aTime, dataMap

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

func parseCSVRecord(record []string, dataMap map[string]map[string]*domain.CovidMetric, categories *utility.Stack, schools *utility.Stack) {
	if record == nil || len(record) <= 0 {
		return
	}

	if utility.CountNonEmptyValuesIn(record) == 2 {
		var category string = record[1]
		_, ok := dataMap[category]
		if !ok {
			dataMap[category] = make(map[string]*domain.CovidMetric)
			for _, aSchool := range *schools {
				//fmt.Println(aSchool)
				var metric domain.CovidMetric
				dataMap[category][aSchool] = &metric
			}
		}
		//fmt.Printf("Found category %s\n", category)
		categories.Push(category)
	} else {
		category, _ := categories.Peak()
		var metricName string = record[1]
		//fmt.Printf("Category %s Metric: %s ... Record %s\n", category, metricName, record)
		if metricName == "Active Cases" {
			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
				metric.ActiveCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)

		} else if metricName == "Total Positive Cases" {
			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
				metric.TotalPositiveCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		} else if metricName == "Probable Cases" {
			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
				metric.ProbableCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		} else if metricName == "Resolved (no longer active)" {
			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
				metric.ResolvedCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		}

	}

}

func updateSchoolMetrics(record []string, category string, dataMap map[string]map[string]*domain.CovidMetric, schools *utility.Stack, metricAssignmentFunction func(*domain.CovidMetric, int)) {
	i := 0
	for _, aSchool := range *schools {
		metric, found := dataMap[category][aSchool]
		if !found {
			log.Fatal("poop")
		}
		recordValue := record[2+i]
		if recordValue != "" {
			value, err := strconv.Atoi(record[2+i])
			if err != nil {
				log.Fatal("Encountered a non numeric value in csv data position")
			}
			metricAssignmentFunction(metric, value)
		}

		i++
	}
}

/*
Store data in redis

organization {
code, description}

category {
code, description}


organizations: [ hhs, hes, hms ...]
categories: [ staff, students ]





	staff
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)



Category
	School
		Date
			Metrics
*/
