package main

import (
	"github.com/pburskey/covid/internal/dao/mysql"
	dao "github.com/pburskey/covid/internal/dao/redis"
	"github.com/pburskey/covid/internal/parser"
	redis_utility "github.com/pburskey/covid/internal/redis"
	"github.com/pburskey/covid/internal/utility"
	"log"
	"os"
)

var daoImpl *dao.DAO

func main() {

	config := utility.LoadConfiguration()

	arguments := os.Args[1:]
	var sourceDirectory string = arguments[0] //"sample_data.csv"
	var targetDirectory string = arguments[1] //"sample_data.csv"

	mysqlconfiguration, err := mysql.Configure(config.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlconfiguration.Close()

	redisConnection := redis_utility.Factory(config.Redis)
	redisDAOImpl := dao.Build(redisConnection)

	covidDB := mysql.Build(mysqlconfiguration, redisDAOImpl)

	shelves := make([]parser.ShelfI, 0)
	shelves = append(shelves, parser.BuildPrettyCSVShelf(covidDB, targetDirectory))
	aParser := parser.BuildRawParser(covidDB, shelves)
	aParser.Parse(sourceDirectory)

}
