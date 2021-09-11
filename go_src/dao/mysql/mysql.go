package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type MySQLConfiguration struct {
	Password string
	UserId   string
	Hostname string
	DbName   string

	dataSource *sql.DB
}

func (me MySQLConfiguration) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", me.UserId, me.Password, me.Hostname, me.DbName)
}

func (me MySQLConfiguration) Close() {
	me.dataSource.Close()
}

func Configure(configurationMap interface{}) (*MySQLConfiguration, error) {

	var sqlConfiguration *MySQLConfiguration

	jsonbody, err := json.Marshal(configurationMap)
	if err != nil {
		// do error check
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(jsonbody, &sqlConfiguration); err != nil {
		// do error check
		fmt.Println(err)
		return nil, err
	}
	dsn := sqlConfiguration.DSN()

	sqlConfiguration.dataSource, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	//
	//ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancelfunc()
	//res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	//if err != nil {
	//	log.Printf("Error %s when creating DB\n", err)
	//	return
	//}
	//no, err := res.RowsAffected()
	//if err != nil {
	//	log.Printf("Error %s when fetching rows", err)
	//	return
	//}
	//log.Printf("rows affected %d\n", no)
	//
	//db.Close()
	//db, err = sql.Open("mysql", dsn(dbname))
	//if err != nil {
	//	log.Printf("Error %s when opening DB", err)
	//	return
	//}
	//defer db.Close()

	sqlConfiguration.dataSource.SetMaxOpenConns(20)
	sqlConfiguration.dataSource.SetMaxIdleConns(20)
	sqlConfiguration.dataSource.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = sqlConfiguration.dataSource.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", sqlConfiguration.DbName)

	return sqlConfiguration, nil

}
