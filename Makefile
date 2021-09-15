APP_NAME=HASD_COVID
export ROOT_DIR=${PWD}
.help:                             ##this help
	@echo
	@echo "Choose a command to run in application: $(APP_NAME)"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo

docker_clean:                       ##cleans up any existing docker instances that this process requires
	./clean_docker.sh redis

start: docker_clean tools           ##Starting this process along with necessary dependencies
	@echo "starting"
	./start.sh &

get_data:                           ##Initiates a data pull request from HASD
	@echo "Pulling data from hasd"
	./doIt.sh

parse_raw:                           ##Initiates a data pull request from HASD
	@echo "Parsing raw data into a more friendly format"
	@echo ${ROOT_DIR}
	go ${ROOT_DIR}/cmd/parse/raw/main.go ${ROOT_DIR}/data/raw ${ROOT_DIR}/data/parsed


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
