package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pburskey/hasd_covid/dao"
	"github.com/pburskey/hasd_covid/domain"
	"log"
	"strconv"
)

type db struct {
	configuration *MySQLConfiguration
	cache         dao.DAO
}

func (me db) GetSchools() ([]*domain.Code, error) {
	var codes []*domain.Code
	var err error
	if codes, err = me.cache.GetSchools(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if codes == nil || len(codes) == 0 {
		res, err := me.configuration.dataSource.Query("SELECT * FROM school")
		defer res.Close()

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		for res.Next() {

			code := &domain.Code{
				Type: domain.SCHOOL,
			}
			err := res.Scan(&code.Id, &code.Description)

			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			codes = append(codes, code)
		}
	}

	return codes, nil
}

func (me db) GetMetric(metricSkey uint) (*domain.CovidMetric, error) {
	var anObject *domain.CovidMetric

	res, err := me.configuration.dataSource.Query("SELECT metric_skey, category_skey, school_skey, ts, active_cases, total_positive_cases, total_probable_cases, resolved FROM metric where metric_skey = ?", metricSkey)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if res.Next() {
		anObject, err = me.metricRowMapper(res)
	}

	return anObject, nil
}

func (me db) metricRowMapper(res *sql.Rows) (*domain.CovidMetric, error) {

	anObject := &domain.CovidMetric{}
	err := res.Scan(&anObject.Id, &anObject.CategoryId, &anObject.SchoolId, &anObject.DateTime, &anObject.ActiveCases, &anObject.TotalPositiveCases, &anObject.ProbableCases, &anObject.ResolvedCases)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return anObject, nil
}

func (me db) GetMetricsBySchool(school *domain.Code) ([]*domain.CovidMetric, error) {
	var metrics []*domain.CovidMetric
	res, err := me.configuration.dataSource.Query("SELECT * FROM metric where school_skey = ?", school.Id)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for res.Next() {
		aMetric, err := me.metricRowMapper(res)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		metrics = append(metrics, aMetric)

	}
	return metrics, nil
}

func (me db) SaveMetric(aMetric *domain.CovidMetric) (*domain.CovidMetric, error) {
	/*
		TABLE IF NOT EXISTS `hasd_covid`.`metric` (
		    `metric_skey` INT NOT NULL AUTO_INCREMENT,
		    `category_skey` INT NOT NULL ,
		    `school_skey` INT NOT NULL ,
		    `ts` TIMESTAMP NOT NULL ,
		    `active_cases` int NOT NULL default 0,
		    `total_positive_cases` int NOT NULL default 0,
		    `total_probable_cases` int NOT NULL default 0,
		    `resolved` int NOT NULL default 0,
	*/
	sql := "INSERT INTO metric(category_skey, school_skey, ts, active_cases, total_positive_cases, total_probable_cases, resolved) VALUES (?,?,?,?,?,?,?)"
	res, err := me.configuration.dataSource.Exec(sql, aMetric.CategoryId, aMetric.SchoolId, aMetric.DateTime, aMetric.ActiveCases, aMetric.TotalPositiveCases, aMetric.ProbableCases, aMetric.ResolvedCases)

	if err != nil {
		panic(err.Error())
	}

	var anID int64
	anID, err = res.LastInsertId()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	aMetric.Id = strconv.FormatInt(anID, 10)
	return aMetric, nil
}

func (me db) SaveSchool(aCode *domain.Code) (*domain.Code, error) {
	sql := "INSERT INTO school(description) VALUES (?)"
	res, err := me.configuration.dataSource.Exec(sql, aCode.Description)

	if err != nil {
		panic(err.Error())
	}

	var anID int64
	anID, err = res.LastInsertId()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	aCode.Id = strconv.FormatInt(anID, 10)
	me.cache.SaveSchool(aCode)
	return aCode, nil
}

func (me db) GetCategories() ([]*domain.Code, error) {
	var codes []*domain.Code
	var err error
	if codes, err = me.cache.GetCategories(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if codes == nil || len(codes) == 0 {
		res, err := me.configuration.dataSource.Query("SELECT * FROM category")
		defer res.Close()

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		for res.Next() {

			code := &domain.Code{
				Type: domain.CATEGORY,
			}
			err := res.Scan(&code.Id, &code.Description)

			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			codes = append(codes, code)
		}
	}

	return codes, nil
}

func (me db) SaveCategory(aCode *domain.Code) (*domain.Code, error) {
	sql := "INSERT INTO category(description) VALUES (?)"
	res, err := me.configuration.dataSource.Exec(sql, aCode.Description)

	if err != nil {
		panic(err.Error())
	}

	var anID int64
	anID, err = res.LastInsertId()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	aCode.Id = strconv.FormatInt(anID, 10)
	me.cache.SaveCategory(aCode)
	return aCode, nil
}

func Build(aDB *MySQLConfiguration, aCache dao.DAO) dao.DAO {
	db := &db{
		configuration: aDB,
		cache:         aCache,
	}
	return db
}
