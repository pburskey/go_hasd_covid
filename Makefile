APP_NAME=HASD_COVID
THIS_FILE := $(lastword $(MAKEFILE_LIST))
export ROOT_DIR=${PWD}
#.help:                             ##this help
#	@echo
#	@echo "Choose a command to run in application: $(APP_NAME)"
#	@echo
#	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
#	@echo

docker_clean:                       ##cleans up any existing docker instances that this process requires
	./clean_docker.sh redis
#
#start: docker_clean tools           ##Starting this process along with necessary dependencies
#	@echo "starting"
#	./start.sh &

get_data:                           ##Initiates a data pull request from HASD
	@echo "Pulling data from hasd"
	./doIt.sh

parse_raw:                           ##Initiates a data pull request from HASD
	rm ${ROOT_DIR}/parse_raw
	@echo "Parsing raw data into a more friendly format"
	@echo ${ROOT_DIR}
	/usr/local/go/bin/go build -o ${ROOT_DIR}/parse_raw ${ROOT_DIR}/cmd/parse/raw
	./parse_raw ${ROOT_DIR} ${ROOT_DIR}/data/raw ${ROOT_DIR}/data/parsed

process_friendly:                           ##Initiates a data pull request from HASD
	#rm ${ROOT_DIR}/process_friendly
	@echo "Parsing raw data into a more friendly format"
	@echo ${ROOT_DIR}
	/usr/local/go/bin/go build -o ${ROOT_DIR}/process_friendly ${ROOT_DIR}/cmd/parse/friendly
	./process_friendly ${ROOT_DIR} ${ROOT_DIR}/data/parsed



report: start                       ##Starts this process, runs a report and shuts down.
	@echo "Running report"

tools:                              ##Performs go get to update needed dependencies
	go get -u github.com/gomodule/redigo
	go get github.com/go-redis/redis
	go get github.com/jinzhu/gorm
	go get github.com/go-sql-driver/mysql
	go get github.com/jinzhu/gorm/dialects/mysql
	go get github.com/rs/cors
	go get github.com/go-sql-driver/mysql
	go get github.com/gorilla/mux
	go get github.com/olekukonko/tablewriter


#
#.PHONY: help build up start down destroy stop restart logs logs-api ps login-timescale login-api db-shell
#help:
#        make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
build:
	docker-compose -f docker-compose.yml build $(c)
up:
	docker-compose -f docker-compose.yml up -d $(c)
start:
	docker-compose -f docker-compose.yml start $(c)
down:
	docker-compose -f docker-compose.yml down $(c)
destroy:
	docker-compose -f docker-compose.yml down -v $(c)
stop:
	docker-compose -f docker-compose.yml stop $(c)
restart:
	docker-compose -f docker-compose.yml stop $(c)
	docker-compose -f docker-compose.yml up -d $(c)
logs:
	docker-compose -f docker-compose.yml logs --tail=100 -f $(c)
logs-api:
	docker-compose -f docker-compose.yml logs --tail=100 -f api
ps:
	docker-compose -f docker-compose.yml ps
login-timescale:
	docker-compose -f docker-compose.yml exec timescale /bin/bash
login-api:
	docker-compose -f docker-compose.yml exec api /bin/bash
db-shell:
	docker-compose -f docker-compose.yml exec timescale psql -Upostgres