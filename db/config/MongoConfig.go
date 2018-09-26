package config

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	logger "github.com/adcamie/adserver/logger"
	"gopkg.in/mgo.v2"
)

var MongoSession *mgo.Session
var sessionPool = map[int]*mgo.Session{}
var mongoBackupPool = map[int]*mgo.Session{}

func InitializeMongo() bool {
	var err error
	fmt.Printf("Creating Session GCP for mongo\n")
	MongoSession, err = mgo.Dial("0.0.0.0:27017,0.0.0.0:27017")
	MongoSession.SetMode(mgo.Eventual, false)
	MongoSession.SetSocketTimeout(5 * time.Minute)
	if err != nil {
		log.Println("Not able to connect to mongo", err)
		logger.ErrorLogger(err.Error(), "Mongo", "MongoSession Creation failed")
		return false
	}
	log.Println("Connection successful at", MongoSession)
	return true
}

func InitializeMongoBackup() {

	fmt.Printf("Creating Backup Session for mongo\n")
	//create 10 sessions
	for i := 0; i < 5; i++ {
		MongoBackup, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    []string{"10.0.0.0:27017"},
			Username: "adcamie",
			Password: "gs#adcamie2017@nov",
			Database: "admin",
		})
		MongoBackup.SetSocketTimeout(5 * time.Minute)
		if err != nil {
			log.Println("Not able to connect to mongo", err)
			logger.ErrorLogger(err.Error(), "MongoBackup", "MongoBackup Creation failed")
		}
		mongoBackupPool[i] = MongoBackup
	}

	log.Println("Mongo Backup Connection successful.")
}

func InitializeMongoSessionPool() bool {
	//create 10 sessions & add to pool
	for i := 0; i < 10; i++ {
		mongo, err := mgo.Dial("0.0.0.0:27017,0.0.0.0:27017")
		mongo.SetMode(mgo.Eventual, false)
		mongo.SetSocketTimeout(5 * time.Minute)
		if err != nil {
			log.Println("Not able to connect to mongo", err)
			logger.ErrorLogger(err.Error(), "MongoBackup", "MongoBackup Creation failed")
			return false
		}
		sessionPool[i] = mongo
	}
	return true
}

//Session Pools
func GetMongoSession() *mgo.Session {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(9)
	log.Print(i)
	return sessionPool[i]
}

func GetMongoBackupSession() *mgo.Session {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(4)
	log.Print(i)
	return mongoBackupPool[i]
}

//close sessions
func ShutdownMongo() {
	MongoSession.Close()
}
func ShutdownMongoSessionPool() {
	for i := 0; i < 10; i++ {
		sessionPool[i].Close()
	}
}

func ShutdownMongoBackup() {
	for i := 0; i < 5; i++ {
		mongoBackupPool[i].Close()
	}
}
