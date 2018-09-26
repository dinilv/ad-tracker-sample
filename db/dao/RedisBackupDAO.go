package dao

/*
As name suggests this class provides functions that is exclusively for reading postback instance critical data.
Contains Writing & Removal of data to different keys.
It validates on each postback try in delayed_job, retry_job and postback_api
*/
import (
	"fmt"

	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	logger "github.com/adcamie/adserver/logger"
)

//validations
func ValidateTransactionIDOnBackup(transactionID string) bool {
	result := db.RedisBackupClient.HGet(constants.Transactions, transactionID).Val()
	if len(result) == 0 {
		fmt.Println("TransactionID is not present in transaction ID Hash")
		return false
	}
	return true
}

//validations
func ValidateSubscriberTransactionIDOnBackup(transactionID string) bool {
	result := db.RedisBackupClient.HGet(constants.PostbackTransactions, transactionID).Val()
	if len(result) == 0 {
		fmt.Println("TransactionID is not present in postback transaction ID Hash")
		return false
	}
	return true
}

func SavePostbackTransaction(transactionID string) {
	pipeline := db.RedisBackupClient.Pipeline()
	pipeline.HSet(constants.Transactions, transactionID, constants.Zero)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving unsent postbacks :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisBackup", "Save tranxn for postback validation")
	}
	defer pipeline.Close()
}

func SavePostbackTransactionOnSubscriber(transactionID string) {

	fmt.Println("in redis saving postbacks", transactionID)

	pipeline := db.RedisBackupClient.Pipeline()
	pipeline.HSet(constants.PostbackTransactions, transactionID, constants.Zero)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving postbacks :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisBackup", "Save tranxn for postback validation")
	}
	defer pipeline.Close()
}

//validations
func ValidateDelayedTransactionIDOnBackup(transactionID string) bool {
	result := db.RedisBackupClient.HGet(constants.DelayedTransactions, transactionID).Val()
	if len(result) == 0 {
		fmt.Println("Delayed TransactionID is not present in transaction ID Hash")
		return false
	}
	return true
}

func SaveDelayedPostbackTransaction(transactionID string) {
	pipeline := db.RedisBackupClient.Pipeline()
	pipeline.HSet(constants.DelayedTransactions, transactionID, constants.Zero)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving unsent postbacks :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisBackup", "Save delayed tranxn for postback validation")
	}
	defer pipeline.Close()
}

//validations
func ValidateRetryTransactionIDOnBackup(transactionID string) bool {
	result := db.RedisBackupClient.HGet(constants.RetryTransactions, transactionID).Val()
	if len(result) == 0 {
		fmt.Println("Retry TransactionID is not present in transaction ID Hash")
		return false
	}
	return true
}

func SaveRetryPostbackTransaction(transactionID string) {
	pipeline := db.RedisBackupClient.Pipeline()
	pipeline.HSet(constants.RetryTransactions, transactionID, constants.Zero)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving unsent postbacks :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisBackup", "Save delayed tranxn for postback validation")
	}
	defer pipeline.Close()
}
