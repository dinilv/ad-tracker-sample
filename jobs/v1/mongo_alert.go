package v1

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

func MongoMasterHealthCheck() bool {
	sess, err := mgo.Dial("10.148.0.2:27017,10.148.0.4:27017")
	sess.SetMode(mgo.Eventual, false)
	if err != nil {
		log.Println("First step")
		fmt.Println(err)
		return true
	}
	err = sess.Ping()
	if err != nil {
		fmt.Println(err)
		return true
	}
	log.Println("MongoDB server is healthy.")
	sess.Close()
	return false
}

func MongoMasterHealthCheckLocal() bool {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("First step")
		fmt.Println(err)
		return true
	}
	err = sess.Ping()
	if err != nil {
		fmt.Println(err)
		return true
	}
	log.Println("MongoDB server is healthy.")
	sess.Close()
	return false
}

func MongoBackupHealthCheck() bool {
	sess, err := mgo.Dial("10.140.0.3:27017")
	if err != nil {
		fmt.Println(err)
		return true
	}
	err = sess.Ping()
	if err != nil {
		fmt.Println(err)
		return true

	}
	log.Println("MongoDB server is healthy.")
	sess.Close()
	return false
}
