package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
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
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func GetRedisConnection() redis.Conn {

	if pool == nil {
		pool = newPool()
		conn := GetRedisConnection()
		defer conn.Close()
		err := Ping(conn)
		if err != nil {
			fmt.Println(err)
		}
	}
	return pool.Get()
}

func CacheKeyExists(aKey string) (exists bool) {
	conn := GetRedisConnection()
	defer conn.Close()
	var aValue int64

	aValue, err := redis.Int64(conn.Do("EXISTS", aKey))

	if err != nil {
		fmt.Println(err)
		return
	}
	exists = (aValue != 0)
	return
}
