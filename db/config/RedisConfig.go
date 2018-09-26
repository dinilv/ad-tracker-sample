package config

import (
	"fmt"

	logger "github.com/adcamie/adserver/logger"
	"github.com/go-redis/redis"
)

var RedisMasterClient *redis.Client
var RedisSlaveClient *redis.Client
var RedisTranxnClient *redis.Client
var RedisBackupClient *redis.Client

//For master redis connection
func InitializeRedisMaster(poolSize int) {
	fmt.Println("Intializing redis-config")
	RedisMasterClient = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6378",
		Password: ""
		DB:       0,
		PoolSize: poolSize,
	})
	err := RedisMasterClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Client Creation Error")
	}
}

//For slave redis connection
func InitializeRedisSlave(poolSize int) {
	fmt.Println("Intializing redis-config")
	RedisSlaveClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		PoolSize: poolSize,
		//Password: "server@123",
	})
	err := RedisSlaveClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisSlave", "Client Creation Error")
	}
}

//For transaction redis connection
func InitializeRedisTranxn(poolSize int) {
	fmt.Println("Intializing redis-config")

	RedisTranxnClient = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "0.0.0.0",
		DB:       0,
		PoolSize: poolSize,
	})
	err := RedisTranxnClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Client Creation Error")
	}
}

//For backup redis connection
func InitializeRedisBackup(poolSize int) {
	fmt.Println("Intializing redis-backup-config")

	RedisBackupClient = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "0.0.0.0",
		DB:       0,
		PoolSize: poolSize,
	})
	err := RedisBackupClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisBackup", "Client Creation Error")
	}
}
