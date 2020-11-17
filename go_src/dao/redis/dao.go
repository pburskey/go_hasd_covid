package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pburskey/hasd_covid/domain"
	redis_utility "github.com/pburskey/hasd_covid/redis"
	"github.com/pburskey/hasd_covid/utility"
	"log"
	"time"
)

func GetMetricsBySchool(aString string) []string {

	var data []string

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	key := fmt.Sprintf("SCHOOL_%s_DATA", aString)
	aValues, err := redis.Values(conn.Do("SMEMBERS", key))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func GetMetricsBySchoolAndCategory(aSchool string, aCategory string) []string {

	var data []string

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	key := fmt.Sprintf("SCHOOL_%s_CATEGORY_%s_DATA", aSchool, aCategory)
	aValues, err := redis.Values(conn.Do("SMEMBERS", key))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func GetMetricsByDate(aString string) []string {

	var data []string

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	key := fmt.Sprintf("DATE_%s_DATA", aString)
	aValues, err := redis.Values(conn.Do("SMEMBERS", key))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func GetMetricsByCategory(aString string) []string {

	var data []string

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	key := fmt.Sprintf("CATEGORY_%s_DATA", aString)
	aValues, err := redis.Values(conn.Do("SMEMBERS", key))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func setMetricInCache(c redis.Conn, metric *domain.DataPoint, key string) error {

	// serialize User object to JSON
	json, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", key, json)
	if err != nil {
		return err
	}

	return nil
}

func GetMetricInCache(key string) (err error, metric *domain.DataPoint) {

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		log.Fatal("Metric does not exist")
	} else if err != nil {
		return err, metric
	}

	metric = &domain.DataPoint{}
	err = json.Unmarshal([]byte(s), metric)

	return err, metric
}

func SaveMetric(category string, school string, dateAndTime time.Time, metric *domain.CovidMetric, identityCounter *utility.Counter) {
	//fmt.Printf("Category %s School: %s ... Date: %s ... Metric %s\n", category, school, dateAndTime, metric)

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	var identity = identityCounter.NextId()
	var key = fmt.Sprintf("METRIC:%d", identity)

	aDateAsString := utility.AsYYYYMMDDHH24MiSS(dateAndTime)

	metricDataPoint := &domain.DataPoint{
		Metric:   *metric,
		Category: category,
		School:   school,
		DateTime: dateAndTime,
		Id:       identity,
	}
	setMetricInCache(conn, metricDataPoint, key)

	////log.Println(key)
	//aValue, error := conn.Do("HSET", key, "category", category, "school", school, "dateTime", aDateAsString, "ActiveCases", metric.ActiveCases, "TotalPositiveCases", metric.TotalPositiveCases, "ProbableCases", metric.ProbableCases, "ResolvedCases", metric.ResolvedCases)
	//if error != nil {
	//	log.Fatal("Unable to store data in redis.... Key:%s Message:%s Error:%s", key, aValue, error)
	//}
	//aValue, error = conn.Do("LPUSH", "METRICS", identity)
	//if error != nil {
	//	log.Fatal("Unable to store data in redis.... Key:%s Message:%s Error:%s", key, aValue, error)
	//}

	//aValue, error = conn.Do("HGETALL", key)

	//aValue, error = conn.Do("LPUSH", "METRICS", fmt.Sprintf("%d", identity))

	_, error := conn.Do("sadd", "CATEGORIES", category)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}
	_, error = conn.Do("sadd", "SCHOOLS", school)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}
	_, error = conn.Do("sadd", "DATES", aDateAsString)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}
	/*
		school data
	*/
	_, error = conn.Do("sadd", fmt.Sprintf("SCHOOL_%s_DATA", school), key)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}
	/*
		category data
	*/
	_, error = conn.Do("sadd", fmt.Sprintf("CATEGORY_%s_DATA", category), key)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}

	/*
		school and category data
	*/
	_, error = conn.Do("sadd", fmt.Sprintf("SCHOOL_%s_CATEGORY_%s_DATA", school, category), key)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}

	/*
		data by date
	*/
	_, error = conn.Do("sadd", fmt.Sprintf("DATE_%s_DATA", aDateAsString), key)
	if error != nil {
		log.Fatal("Unable to store data in redis.... ", error)
	}
	//log.Println(error)
	//log.Println(aValue)

}

func GetDates() (data []string, err error) {

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "DATES"))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return
}

func GetCategories() (data []string, err error) {

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "CATEGORIES"))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return
}

func GetSchools() (data []string, err error) {

	conn := redis_utility.GetRedisConnection()
	defer conn.Close()

	aValues, err := redis.Values(conn.Do("SMEMBERS", "SCHOOLS"))

	if err != nil {
		log.Fatal(err)
	}

	if err := redis.ScanSlice(aValues, &data); err != nil {
		log.Fatal(err)
	}
	return
}
