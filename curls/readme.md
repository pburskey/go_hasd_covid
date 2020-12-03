curl --request GET localhost:8080/api/v1/schools
curl --request GET localhost:8080/api/v1/categories
curl --request GET localhost:8080/api/v1/dates

curl --request GET localhost:8080/api/v1/category/Staff/metrics
curl --request GET localhost:8080/api/v1/school/HES/metrics
curl --request GET localhost:8080/api/v1/date/20201021104718/metrics

curl --request GET localhost:8080/api/v1/metric/METRIC:313

http://127.0.0.1:8080/api/v1/school/%s/category/%s/metrics
curl --request GET localhost:8080/api/v1/school/HES/category/Students/metrics


curl -ivk \
localhost:8080/api/v1/school/HES/category/Students/metrics \
-H "Content-Type: application/json" \
-H "resolve-metric-detail: 1" \
-X GET



curl -ivk \
localhost:8080/api/v1/school/HES/category/Students/metricDetails \
-H "Content-Type: application/json" \
-X GET