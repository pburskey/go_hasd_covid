package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gomodule/redigo/redis"
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

	aTime, err := parseYYYYMMDDHH24MiSS(dateAndTimePortion)
	if err != nil {
		log.Fatal("Unable to parse date and time portion from file name: %s", fileName)
	}
	return aTime
}

func parseCSV(fileName string) (time.Time, map[string]map[string]*CovidMetric) {

	aTime := parseDateFromFileName(fileName)
	//fmt.Println("Date and time: %s from file name: %s", aTime, fileName)
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file: %s for processing", fileName, err)
	}
	r := csv.NewReader(csvFile)
	if r == nil {
		return aTime, nil
	}

	dataMap := make(map[string]map[string]*CovidMetric)
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
	//log.Println(dataMap)
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

type CovidMetric struct {
	ActiveCases        int
	TotalPositiveCases int
	ProbableCases      int
	ResolvedCases      int
}

type DataPoint struct {
	CovidMetric
	category string
	school   string
	dateTime string
}

func parseCSVRecord(record []string, dataMap map[string]map[string]*CovidMetric, categories *Stack, schools *Stack) {
	if record == nil || len(record) <= 0 {
		return
	}

	if countNonEmptyValuesIn(record) == 2 {
		var category string = record[1]
		_, ok := dataMap[category]
		if !ok {
			dataMap[category] = make(map[string]*CovidMetric)
			for _, aSchool := range *schools {
				//fmt.Println(aSchool)
				var metric CovidMetric
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
			metricAssignmentFunction := func(metric *CovidMetric, value int) {
				metric.ActiveCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)

		} else if metricName == "Total Positive Cases" {
			metricAssignmentFunction := func(metric *CovidMetric, value int) {
				metric.TotalPositiveCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		} else if metricName == "Probable Cases" {
			metricAssignmentFunction := func(metric *CovidMetric, value int) {
				metric.ProbableCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		} else if metricName == "Resolved (no longer active)" {
			metricAssignmentFunction := func(metric *CovidMetric, value int) {
				metric.ResolvedCases = value
			}
			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
		}

	}

}

func updateSchoolMetrics(record []string, category string, dataMap map[string]map[string]*CovidMetric, schools *Stack, metricAssignmentFunction func(*CovidMetric, int)) {
	i := 0
	for _, aSchool := range *schools {
		metric, found := dataMap[category][aSchool]
		if !found {
			log.Fatal("poop")
		}
		value, err := strconv.Atoi(record[2+i])
		if err != nil {
			log.Fatal("Encountered a non numeric value in csv data position")
		}
		metricAssignmentFunction(metric, value)

		i++
	}
}

func saveMetric(category string, school string, dateAndTime time.Time, metric *CovidMetric) {
	//fmt.Printf("Category %s School: %s ... Date: %s ... Metric %s\n", category, school, dateAndTime, metric)

	conn := getRedisConnection()
	defer conn.Close()

	var identity = identityCounter.nextId()
	var key = fmt.Sprintf("METRIC:%d", identity)

	aDateAsString := asYYYYMMDDHH24MiSS(dateAndTime)

	//log.Println(key)
	aValue, error := conn.Do("HSET", key, "category", category, "school", school, "dateTime", aDateAsString, "ActiveCases", metric.ActiveCases, "TotalPositiveCases", metric.TotalPositiveCases, "ProbableCases", metric.ProbableCases, "ResolvedCases", metric.ResolvedCases)
	if error != nil {
		log.Fatal("Unable to store data in redis.... Key:%s Message:%s Error:%s", key, aValue, error)
	}
	//aValue, error = conn.Do("LPUSH", "METRICS", identity)
	//if error != nil {
	//	log.Fatal("Unable to store data in redis.... Key:%s Message:%s Error:%s", key, aValue, error)
	//}

	//aValue, error = conn.Do("HGETALL", key)

	//aValue, error = conn.Do("LPUSH", "METRICS", fmt.Sprintf("%d", identity))

	aValue, error = conn.Do("sadd", "CATEGORIES", category)
	aValue, error = conn.Do("sadd", "SCHOOLS", school)
	aValue, error = conn.Do("sadd", "DATES", aDateAsString)

	/*
		school data
	*/
	aValue, error = conn.Do("sadd", fmt.Sprintf("SCHOOL_%s_DATA", school), key)

	/*
		category data
	*/
	aValue, error = conn.Do("sadd", fmt.Sprintf("CATEGORY_%s_DATA", category), key)

	/*
		data by date
	*/
	aValue, error = conn.Do("sadd", fmt.Sprintf("DATE_%s_DATA", aDateAsString), key)

	log.Println(error)
	log.Println(aValue)

}

// set executes the redis SET command
func set(c redis.Conn) error {
	_, err := c.Do("SET", "Favorite Movie", "Repo Man")
	if err != nil {
		return err
	}
	_, err = c.Do("SET", "Release Year", 1984)
	if err != nil {
		return err
	}
	return nil
}

// get executes the redis GET command
func get(c redis.Conn) error {

	// Simple GET example with String helper
	key := "Favorite Movie"
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	// Simple GET example with Int helper
	key = "Release Year"
	i, err := redis.Int(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %d\n", key, i)

	// Example where GET returns no results
	key = "Nonexistent Key"
	s, err = redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		fmt.Printf("%s does not exist\n", key)
	} else if err != nil {
		return err
	} else {
		fmt.Printf("%s = %s\n", key, s)
	}

	return nil
}

//
//func setStruct(c redis.Conn) error {
//
//	const objectPrefix string = "user:"
//
//	usr := User{
//		Username:  "otto",
//		MobileID:  1234567890,
//		Email:     "ottoM@repoman.com",
//		FirstName: "Otto",
//		LastName:  "Maddox",
//	}
//
//	// serialize User object to JSON
//	json, err := json.Marshal(usr)
//	if err != nil {
//		return err
//	}
//
//	// SET object
//	_, err = c.Do("SET", objectPrefix+usr.Username, json)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getStruct(c redis.Conn) error {
//
//	const objectPrefix string = "user:"
//
//	username := "otto"
//	s, err := redis.String(c.Do("GET", objectPrefix+username))
//	if err == redis.ErrNil {
//		fmt.Println("User does not exist")
//	} else if err != nil {
//		return err
//	}
//
//	usr := User{}
//	err = json.Unmarshal([]byte(s), &usr)
//
//	fmt.Printf("%+v\n", usr)
//
//	return nil
//
//}
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
