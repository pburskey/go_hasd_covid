package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pburskey/covid/internal/dao"
	"github.com/pburskey/covid/internal/dao/mysql"
	"github.com/pburskey/covid/internal/dao/redis"
	"github.com/pburskey/covid/internal/domain"
	redis_utility "github.com/pburskey/covid/internal/redis"
	"github.com/pburskey/covid/internal/utility"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var daoImpl dao.DAO

func main() {

	arguments := os.Args[1:]
	config := utility.LoadConfiguration(arguments[0])

	mysqlconfiguration, err := mysql.Configure(config.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlconfiguration.Close()

	redisConnection := redis_utility.Factory(config.Redis)
	redisDAOImpl := redis.Build(redisConnection)

	daoImpl = mysql.Build(mysqlconfiguration, redisDAOImpl)

	r := mux.NewRouter()
	r.Use(CORS)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)

	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)
	api.HandleFunc("/schools", schools).Methods(http.MethodGet)
	api.HandleFunc("/categories", categories).Methods(http.MethodGet)
	//api.HandleFunc("/dates", dates).Methods(http.MethodGet)
	api.HandleFunc("/category/{aCategory}/metrics", metricsByCategory).Methods(http.MethodGet)
	api.HandleFunc("/school/{aSchool}/metrics", metricsBySchool).Methods(http.MethodGet)
	api.HandleFunc("/school/{aSchool}/category/{aCategory}/metrics", metricsBySchoolAndCategory).Methods(http.MethodGet)
	api.HandleFunc("/school/{aSchool}/category/{aCategory}/metricDetails", metricDetailsBySchoolAndCategory).Methods(http.MethodGet)

	//api.HandleFunc("/date/{aDate}/metrics", metricsByDate).Methods(http.MethodGet)
	api.HandleFunc("/metric/{aMetric}", metric).Methods(http.MethodGet)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	cors(r)
	//log.Fatal(http.ListenAndServe(":8080", r))
	http.ListenAndServe(":8080", setHeaders(r))

}

//
//func save(aDateAndTime time.Time, dataMap map[string]map[string]*domain.CovidMetric, counter *utility.Counter, daoImpl *dao.DAO) {
//	if dataMap != nil && len(dataMap) > 0 {
//
//		for organization, schools := range dataMap {
//			//log.Println(organization)
//
//			for schoolName, metric := range schools {
//				//log.Println(schoolName)
//				//log.Println(metric)
//				daoImpl.SaveMetric(organization, schoolName, aDateAndTime, metric, counter)
//			}
//
//		}
//	}
//
//}

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

func categories(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//for len(aValues) > 0 {
	//	var aString string
	//	values, err = redis.Scan(values, &aString)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//}

	data, err := daoImpl.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func metricsByCategory(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aCategory"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Category is required"}`))
		return
	}
	categories, _ := daoImpl.GetCategories()
	metricData, err := daoImpl.GetMetricsByCategory(domain.FindCodeByID(categories, aParam))
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

	data, err := daoImpl.GetSchools()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	//aHeaderValue := r.Header.Get("resolve-metric-detail")
	//var resolveMetricDetail bool
	//if len(aHeaderValue) > 0 && aHeaderValue == "1"{
	//	resolveMetricDetail = true
	//}
	aParam, ok := pathParams["aSchool"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required"}`))
		return
	}

	codes, _ := daoImpl.GetSchools()
	code := domain.FindCodeByID(codes, aParam)

	metricData, err := daoImpl.GetMetricsBySchool(code)
	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metricsBySchoolAndCategory(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aSchool, ok := pathParams["aSchool"]

	if !ok || len(aSchool) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required", "message": "Category is required"}`))
		return
	}

	aCategory, ok := pathParams["aCategory"]
	if !ok || len(aCategory) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required", "message": "Category is required"}`))
		return
	}

	schools, _ := daoImpl.GetSchools()
	school := domain.FindCodeByDescription(schools, aSchool)

	categories, _ := daoImpl.GetCategories()
	category := domain.FindCodeByDescription(categories, aCategory)

	metricData, err := daoImpl.GetMetricsBySchoolAndCategory(school, category)
	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func metricDetailsBySchoolAndCategory(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aSchool, ok := pathParams["aSchool"]

	if !ok || len(aSchool) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required", "message": "Category is required"}`))
		return
	}

	aCategory, ok := pathParams["aCategory"]
	if !ok || len(aCategory) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "School is required", "message": "Category is required"}`))
		return
	}

	schools, _ := daoImpl.GetSchools()
	school := domain.FindCodeByID(schools, aSchool)

	categories, _ := daoImpl.GetCategories()
	category := domain.FindCodeByID(categories, aCategory)

	metricData, err := daoImpl.GetMetricsBySchoolAndCategory(school, category)

	sort.Slice(metricData, func(i int, j int) bool {
		return metricData[i].DateTime.Before(metricData[j].DateTime)
	})

	json, err := json.Marshal(&metricData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//
//func dates(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Content-Type", "application/json")
//
//	data, err := daoImpl.GetDates()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	json, err := json.Marshal(&data)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(json)
//}
//
//func metricsByDate(w http.ResponseWriter, r *http.Request) {
//	pathParams := mux.Vars(r)
//	w.Header().Set("Content-Type", "application/json")
//
//	aParam, ok := pathParams["aDate"]
//	if !ok || len(aParam) == 0 {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte(`{"message": "Date as yyyymmddhh24miss is required"}`))
//		return
//	}
//	metricData := daoImpl.GetMetricsByDate(aParam)
//	json, err := json.Marshal(&metricData)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(json)
//}

func metric(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	aParam, ok := pathParams["aMetric"]
	if !ok || len(aParam) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Date as yyyymmddhh24miss is required"}`))
		return
	}
	skey, err := strconv.Atoi(aParam)
	metricData, err := daoImpl.GetMetric(skey)
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
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}

	// process the request...
}

func handler(w http.ResponseWriter, req *http.Request) {
	// ...
	enableCors(&w)
	// ...
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		return
	})
}

func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request (not recommended in production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//only allow GET, POST, and OPTIONS
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//Since I was building a REST API that returned JSON, I set the content type to JSON here.
		w.Header().Set("Content-Type", "application/json")
		//Allow requests to have the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		//if it's just an OPTIONS request, nothing other than the headers in the response is needed.
		//This is essential because you don't need to handle the OPTIONS requests in your handlers now
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
