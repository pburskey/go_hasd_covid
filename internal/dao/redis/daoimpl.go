package redis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/pburskey/covid/internal/dao"
	"github.com/pburskey/covid/internal/domain"
	redis_utility "github.com/pburskey/covid/internal/redis"
	"log"
)

type db struct {
	redis *redis_utility.RedisConnection
}

func (me db) GetMetricsByCategory(code *domain.Code) ([]*domain.CovidMetric, error) {
	panic("implement me")
}

func (me db) GetMetricsBySchoolAndCategory(school *domain.Code, category *domain.Code) ([]*domain.CovidMetric, error) {
	panic("implement me")
}

func (me db) GetSchools() ([]*domain.Code, error) {
	var data []*domain.Code

	conn := me.redis.GetRedisConnection()
	defer conn.Close()

	//key := fmt.Sprintf("SCHOOL_%s_DATA", aString)
	values, err := redis.Values(conn.Do("SMEMBERS", "SCHOOLS"))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, aValue := range values {
		var anObject *domain.Code
		var bytes []byte
		if bytes, err = redis.Bytes(aValue, nil); err != nil {
			log.Fatal(err)
			return nil, err
		}

		if err = json.Unmarshal(bytes, &anObject); err != nil {
			log.Fatal(err)
			return nil, err
		}

		data = append(data, anObject)
	}

	return data, err
}

func (me db) GetCategories() ([]*domain.Code, error) {
	var data []*domain.Code

	conn := me.redis.GetRedisConnection()
	defer conn.Close()

	//key := fmt.Sprintf("SCHOOL_%s_DATA", aString)
	values, err := redis.Values(conn.Do("SMEMBERS", "CATEGORIES"))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, aValue := range values {
		var anObject *domain.Code
		var bytes []byte
		if bytes, err = redis.Bytes(aValue, nil); err != nil {
			log.Fatal(err)
			return nil, err
		}

		if err = json.Unmarshal(bytes, &anObject); err != nil {
			log.Fatal(err)
			return nil, err
		}

		data = append(data, anObject)
	}

	return data, err
}

func (me db) GetMetric(u int) (*domain.CovidMetric, error) {
	panic("implement me")
}

func (me db) GetMetricsBySchool(school *domain.Code) ([]*domain.CovidMetric, error) {
	panic("implement me")
}

func (me db) SaveMetric(metric *domain.CovidMetric) (*domain.CovidMetric, error) {
	panic("implement me")
}

func (me db) SaveSchool(school *domain.Code) (*domain.Code, error) {
	conn := me.redis.GetRedisConnection()
	defer conn.Close()

	bytes, err := json.Marshal(school)

	if _, err = conn.Do("sadd", "SCHOOLS", bytes); err != nil {
		log.Fatal("Unable to store data in redis.... ", err)
		return nil, err
	}

	return school, nil
}

func (me db) SaveCategory(category *domain.Code) (*domain.Code, error) {
	conn := me.redis.GetRedisConnection()
	defer conn.Close()

	bytes, err := json.Marshal(category)

	if _, err = conn.Do("sadd", "CATEGORIES", bytes); err != nil {
		log.Fatal("Unable to store data in redis.... ", err)
		return nil, err
	}

	return category, nil
}

func Build(aRedis *redis_utility.RedisConnection) dao.DAO {
	db := &db{
		redis: aRedis,
	}
	return db
}
