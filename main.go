package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var pool *redis.Pool
var identityCounter Counter

func main() {

	//
	//log.Println(identityCounter.next())
	//log.Println(identityCounter.next())
	//log.Println(identityCounter.next())

	pool = newPool()
	conn := getRedisConnection()
	defer conn.Close()
	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}

	//conn.Send("HMSET", "album:1", "title", "Red", "rating", 5)
	//conn.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
	//conn.Send("HMSET", "album:3", "title", "Beat")
	////conn.Send("LPUSH", "albums", "1")
	////conn.Send("LPUSH", "albums", "2")
	////conn.Send("LPUSH", "albums", "3")
	//values, err := redis.Values(conn.Do("SORT", "albums",
	//	"BY", "album:*->rating",
	//	"GET", "album:*->title",
	//	"GET", "album:*->rating"))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for len(values) > 0 {
	//	var title string
	//	rating := -1 // initialize to illegal value to detect nil.
	//	values, err = redis.Scan(values, &title, &rating)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	if rating == -1 {
	//		fmt.Println(title, "not-rated")
	//	} else {
	//		fmt.Println(title, rating)
	//	}
	//}

	arguments := os.Args[1:]
	var fileOrDirectory string = arguments[0] //"sample_data.csv"

	fileMode, err := os.Stat(fileOrDirectory)
	if err != nil {
		log.Fatal("Encountered error opening: %s", fileOrDirectory, err)
	}
	if fileMode.Mode().IsDir() {
		files, err := ioutil.ReadDir(fileOrDirectory)
		if err != nil {
			log.Fatal("Encountered error reading file names in directory: %s", fileOrDirectory, err)
		}
		for _, aFileName := range files {
			//fmt.Println(aFileName.Name())

			mungedName := (fileOrDirectory + "/" + aFileName.Name())
			save(parseCSV(mungedName))
		}
	} else {
		//log.Print(fileOrDirectory)
		save(parseCSV(fileOrDirectory))
	}

	//
	//letters := [...]string{"a", "b", "c", "d", "", ""}
	//aSlice := letters[4:]
	//fmt.Printf("letters: %s\n", aSlice)
	//fmt.Printf("Only Empties: %s\n", sliceRangeContainsOnlyEmpties(aSlice))
	//fmt.Printf("Has Non Empty: %s\n", sliceRangeContainsNonEmptyValue(aSlice))
	//fmt.Printf("Number of Empties: %s\n", countEmptyValuesIn(aSlice))
	//fmt.Printf("Number of Non Empty values: %s\n", countNonEmptyValuesIn(aSlice))
	//
	//aSlice = letters[0:4]
	//fmt.Printf("letters: %s\n", aSlice)
	//fmt.Printf("Only Empties: %s\n", sliceRangeContainsOnlyEmpties(aSlice))
	//fmt.Printf("Has Non Empty: %s\n", sliceRangeContainsNonEmptyValue(aSlice))
	//fmt.Printf("Number of Empties: %s\n", countEmptyValuesIn(aSlice))
	//fmt.Printf("Number of Non Empty values: %s\n", countNonEmptyValuesIn(aSlice))

	aValue, err := redis.Values(conn.Do("SMEMBERS", "SCHOOLS"))
	log.Println(aValue)

	/*
	   prints out all keys
	*/
	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		// handle error
	}
	for _, key := range keys {
		fmt.Println(key)
	}

	values, err := redis.Values(conn.Do("SORT", "METRICS",
		//"BY", "METRIC:*->category",
		//"BY", "METRIC:*->school",
		//"BY", "METRIC:*->dateTime",
		"GET", "METRIC:->category",
		"GET", "METRIC:->school",
		"GET", "METRIC:->dateTime",
		"GET", "METRIC:->ActiveCases",
		"GET", "METRIC:->TotalPositiveCases",
		"GET", "METRIC:->ProbableCases",
		"GET", "METRIC:->ResolvedCases"))

	if err != nil {
		fmt.Println(err)
		return
	}

	var data []DataPoint
	if err := redis.ScanSlice(values, &data); err != nil {
		fmt.Println(err)
		return
	}
	for len(values) > 0 {
		var aString string
		values, err = redis.Scan(values, &aString)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

func save(aDateAndTime time.Time, dataMap map[string]map[string]*CovidMetric) {
	for organization, schools := range dataMap {
		//log.Println(organization)

		for schoolName, metric := range schools {
			//log.Println(schoolName)
			//log.Println(metric)
			saveMetric(organization, schoolName, aDateAndTime, metric)
		}

	}
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

type HASDCovidMetric struct {
	id          uint
	category    string
	school      string
	dateAndTime time.Time
	metric      *CovidMetric
}

func getRedisConnection() redis.Conn {
	return pool.Get()
}

type Counter struct {
	id uint
}

func (c *Counter) next() *Counter {
	c.id = c.id + 1
	return c
}
func (c *Counter) nextId() uint {
	return (*c).next().id
}
