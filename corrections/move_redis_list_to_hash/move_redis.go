package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	logger "github.com/adcamie/adserver/logger"
	redis "github.com/go-redis/redis"
)

func main() {
	SaveOfferAffiliate()
	db.InitializeMongo()
	//initialize redis
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "10.140.0.2:6378",
		Password: "server@123",
		DB:       0,
		PoolSize: 1000,
	})

	err := RedisClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisClient", "Client Creation Error")
	}

	//query with each pattern on redis live server
	results := RedisClient.SMembers(constants.ExhaustedOfferStack)

	keys := results.Val()

	//exit critieria
	log.Println("Length Of Keys", len(keys))

	//for saving and removing keys
	pipelineRedis := RedisClient.Pipeline()
	for _, key := range keys {
		log.Println("key", key)
		pipelineRedis.HSet(constants.ExhaustedOfferHash, key, constants.Zero)
		log := &model.RotationStack{
			OfferID:   key,
			AddedDate: time.Now().UTC(),
			Event:     "Added by Default Job",
		}
		dao.InsertToMongo(constants.MongoDB, constants.ExhaustedOffer, log)
	}

	_, err = pipelineRedis.Exec()
	if err != nil {
		fmt.Println("Redis error while loading redis-keys :", err.Error())
	}

	//close the pipelines
	defer pipelineRedis.Close()

}

func SaveOfferAffiliate() {
	db.InitializeMongo()
	//initialize redis
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "10.140.0.2:6378",
		Password: "server@123",
		DB:       0,
		PoolSize: 1000,
	})

	err := RedisClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
		go logger.ErrorLogger(err.Error(), "RedisClient", "Client Creation Error")
	}

	//query with each pattern on redis live server
	results := RedisClient.SMembers(constants.ExhaustedOfferAffiliateStack)

	keys := results.Val()

	//exit critieria
	log.Println("Length Of Keys", len(keys))

	//for saving and removing keys
	pipelineRedis := RedisClient.Pipeline()
	for _, key := range keys {
		log.Println("key", key)
		pipelineRedis.HSet(constants.ExhaustedOfferAffiliateHash, key, constants.Zero)
		splittedKeys := strings.Split(key, constants.Separator)
		if len(splittedKeys) > 1 {
			log := &model.RotationStack{
				OfferID:     splittedKeys[0],
				AffiliateID: splittedKeys[1],
				AddedDate:   time.Now().UTC(),
				Event:       "Added by Default Job",
			}
			dao.InsertToMongo(constants.MongoDB, constants.ExhaustedOfferAffiliate, log)
		}

	}

	_, err = pipelineRedis.Exec()
	if err != nil {
		fmt.Println("Redis error while loading redis-keys :", err.Error())
	}

	//close the pipelines
	defer pipelineRedis.Close()

}
