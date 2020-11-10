GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.10.0

go mod init hasd_covid

go get github.com/cucumber/godog/cmd/godog

mkdir features
echo \
'
Feature: Sunny Day
  This is a sunny day example

Scenario: The HASD web site is still available
  Given the HASD has a covid site
  When I consume the url
  Then the http status should be "ok"' >> features/sunnyday.feature

godog

go get github.com/gomodule/redigo
go get github.com/go-redis/redis

go get -u github.com/jinzhu/gorm
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jinzhu/gorm/dialects/mysql
go get -u github.com/rs/cors

go get -u github.com/go-sql-driver/mysql