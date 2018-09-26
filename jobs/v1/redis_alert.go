package v1

import (
	"log"

	"github.com/go-redis/redis"
)

//redis master
func RedisMasterHealthCheck() bool {

	client := redis.NewClient(&redis.Options{
		Addr:     "10.140.0.2:6378",
		Password: "server@123",
		DB:       0,
		PoolSize: 1,
	})
	defer client.Close()
	resp := client.Ping()
	if resp.Err() != nil {
		log.Println(resp.Err())
		return true

	}
	log.Println("RedisDB server is healthy.")
	return false

}
func RedisTranxnHealthCheck() bool {

	client := redis.NewClient(&redis.Options{
		Addr:     "10.140.0.2:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 1,
	})
	defer client.Close()
	resp := client.Ping()
	if resp.Err() != nil {
		log.Println(resp.Err())
		return true

	}
	log.Println("RedisDB server is healthy.")
	return false

}

func RedisBackupHealthCheck() bool {

	client := redis.NewClient(&redis.Options{
		Addr:     "10.140.0.3:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 1,
	})
	defer client.Close()
	resp := client.Ping()
	if resp.Err() != nil {
		log.Println(resp.Err())
		return true

	}
	log.Println("RedisDB server is healthy.")
	return false

}
