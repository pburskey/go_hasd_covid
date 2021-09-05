package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pburskey/hasd_covid/utility"
	"log"
)

type RedisConnection struct {
	pool   *redis.Pool
	config utility.RedisConfiguration
}

func Factory(aconfig utility.RedisConfiguration) *RedisConnection {
	return &RedisConnection{
		config: aconfig,
	}
}

func (r *RedisConnection) newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.

		Dial: func() (redis.Conn, error) {

			var c, err = redis.Dial("tcp", ":6379", redis.DialPassword(r.config.Password))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func Ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil || s != "OK" {
		return err
	}
	//
	//fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func (r *RedisConnection) GetRedisConnection() redis.Conn {

	if r.pool == nil {
		r.pool = r.newPool()
		conn := r.GetRedisConnection()
		defer conn.Close()
		err := Ping(conn)
		if err != nil {
			log.Fatal(err)
		}
	}
	return r.pool.Get()
}

func (r *RedisConnection) CacheKeyExists(aKey string) (exists bool) {
	conn := r.GetRedisConnection()
	defer conn.Close()
	var aValue int64

	aValue, err := redis.Int64(conn.Do("EXISTS", aKey))

	if err != nil {
		log.Fatal(err)
		return
	}
	exists = (aValue != 0)
	return
}
