package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)

	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)
	api.HandleFunc("/schools", schools).Methods(http.MethodGet)
	api.HandleFunc("/categories", categories).Methods(http.MethodGet)
	api.HandleFunc("/dates", dates).Methods(http.MethodGet)
	api.HandleFunc("/category/{aCategory}/metrics", metricsByCategory).Methods(http.MethodGet)
	api.HandleFunc("/school/{aSchool}/metrics", metricsBySchool).Methods(http.MethodGet)
	api.HandleFunc("/date/{aDate}/metrics", metricsByDate).Methods(http.MethodGet)
	api.HandleFunc("/metric/{aMetric}", metric).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func save(aDateAndTime time.Time, dataMap map[string]map[string]*CovidMetric) {
	if dataMap != nil && len(dataMap) > 0 {
		for organization, schools := range dataMap {
			//log.Println(organization)

			for schoolName, metric := range schools {
				//log.Println(schoolName)
				//log.Println(metric)
				saveMetric(organization, schoolName, aDateAndTime, metric)
			}

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

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
}

func cacheKeyExists(aKey string) (exists bool) {
	conn := getRedisConnection()
	defer conn.Close()
	var aValue int64

	aValue, err := redis.Int64(conn.Do("EXISTS", aKey))

	if err != nil {
		fmt.Println(err)
		return
	}
	exists = (aValue != 0)
	return
}

func categories(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	conn := getRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "CATEGORIES"))
	log.Println(aValues)

	if err != nil {
		fmt.Println(err)
		return
	}

	var data []string
	if err := redis.ScanSlice(aValues, &data); err != nil {
		fmt.Println(err)
		return
	}
	//for len(aValues) > 0 {
	//	var aString string
	//	values, err = redis.Scan(values, &aString)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//}

	json, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metricsByCategory(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aCategory"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Category is required"}`))
		return
	}
	metricData := getMetricsByCategory(aParam)
	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func schools(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	conn := getRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "SCHOOLS"))
	log.Println(aValues)

	if err != nil {
		fmt.Println(err)
		return
	}

	var data []string
	if err := redis.ScanSlice(aValues, &data); err != nil {
		fmt.Println(err)
		return
	}

	json, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metricsBySchool(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aSchool"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required"}`))
		return
	}
	metricData := getMetricsBySchool(aParam)
	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func dates(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	conn := getRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "DATES"))
	log.Println(aValues)

	if err != nil {
		fmt.Println(err)
		return
	}

	var data []string
	if err := redis.ScanSlice(aValues, &data); err != nil {
		fmt.Println(err)
		return
	}

	json, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metricsByDate(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aDate"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Date as yyyymmddhh24miss is required"}`))
		return
	}
	metricData := getMetricsByDate(aParam)
	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metric(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aMetric"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Date as yyyymmddhh24miss is required"}`))
		return
	}
	err, metricData := getMetricInCache(aParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
