package dao

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	logger "github.com/adcamie/adserver/logger"
)

//copy details from redis master
func SaveOfferOnTranxn(offerID string, offerType string, template string, group string, advertiser string, countryIds []string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	pipeline.HSet(constants.RedisOfferType, offerID, offerType)
	if len(advertiser) > 0 {
		pipeline.HSet(constants.RedisOfferAdvertiser, offerID, advertiser)
	} else {
		pipeline.HSet(constants.RedisOfferAdvertiser, offerID, "")
	}
	pipeline.HSet(constants.RedisOfferAdvertiserTemplate, offerID, template)
	pipeline.HSet(constants.RedisOfferType, offerID, offerType)
	pipeline.Del(constants.OfferGeoTargetting + constants.Separator + offerID)
	for _, countryID := range countryIds {
		pipeline.SAdd(constants.OfferGeoTargetting+constants.Separator+offerID, countryID)
	}
	pipeline.HSet(constants.RedisOfferGroup, offerID, group)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Offer")
	}
	defer pipeline.Close()
}
func SaveExhaustedOfferOnTranxn(offID string) {
	err := db.RedisTranxnClient.SAdd(constants.ExhaustedOfferStack, offID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save exhausted Offer")
	}
}

func SaveExhaustedOfferAffiliateOnTranxn(offerID string, affiliateID string) {
	err := db.RedisTranxnClient.SAdd(constants.ExhaustedOfferAffiliateStack, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Exhausted Offer Affiliate")
	}
}

func SaveExhaustedOfferMapOnTranxn(offID string) {
	err := db.RedisTranxnClient.HSet(constants.ExhaustedOfferHash, offID, constants.Zero).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save exhausted Offer on Hash")
	}
}

func SaveExhaustedOfferAffiliateMapOnTranxn(offerID string, affiliateID string) {
	err := db.RedisTranxnClient.HSet(constants.ExhaustedOfferAffiliateHash, offerID+constants.Separator+affiliateID, constants.Zero).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save exhausted Offer on Hash")
	}
}

//remove keys
func RemoveOfferInExhaustedOnTranxn(offID string) {
	err := db.RedisTranxnClient.SRem(constants.ExhaustedOfferStack, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove Offer In Exhausted")
	}
}

func RemoveOfferAffiliateInExhaustedOnTranxn(offerID string, affiliateID string) {
	err := db.RedisTranxnClient.SRem(constants.ExhaustedOfferAffiliateStack, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer affiliate in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove Offer Affiliate MQF")
	}
}

func RemoveOfferInExhaustedOnTranxnHash(offID string) {
	err := db.RedisTranxnClient.HDel(constants.ExhaustedOfferHash, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove Offer In Exhausted On Hash")
	}
}

func RemoveOfferAffiliateInExhaustedOnTranxnHash(offerID string, affiliateID string) {
	err := db.RedisTranxnClient.HDel(constants.ExhaustedOfferAffiliateHash, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer affiliate in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove Offer Affiliate MQF on Hash")
	}
}

//details specific to tranxn
func SaveAffiliate(affID string, mqf string, template string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	pipeline.HSet(constants.RedisAffiliateMQF, affID, mqf)
	pipeline.HSet(constants.RedisAffiliateMediaTemplate, affID, template)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Affiliate")
	}
	defer pipeline.Close()
}

func SaveOfferAffiliateMQF(offerID string, affID string, mqf string, template string, rotation bool, resetCounter bool) {

	log.Println("From setting MQF", offerID, "OfferID", affID, "affID")

	pipeline := db.RedisTranxnClient.Pipeline()
	//template
	if len(template) > 0 {
		pipeline.HSet(constants.RedisOfferAffiiatePostback, offerID+constants.Separator+affID, template)
	} else {
		pipeline.HSet(constants.RedisOfferAffiiatePostback, offerID+constants.Separator+affID, "")
	}
	pipeline.HSet(constants.RedisOfferAffiliateMQF, offerID+constants.Separator+affID, mqf)

	//reset rotation status
	if rotation {
		SaveExhaustedOfferAffiliate(offerID, affID)
		SaveExhaustedOfferAffiliateMap(offerID, affID)
	} else {
		RemoveOfferAffiliateInExhausted(offerID, affID)
		RemoveOfferAffiliateInExhaustedHash(offerID, affID)
	}
	//reset counters
	if resetCounter {
		SaveCounters(offerID, affID)
	}
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Offer Affiliate")
	}
	defer pipeline.Close()

}

func SaveOfferAffiliateMQFOnPostback(offerID string, affID string, mqf string) {

	log.Println("From setting MQF on Offer and Media On postback", offerID, "OfferID", affID, "affID")

	pipeline := db.RedisTranxnClient.Pipeline()
	pipeline.HSet(constants.RedisOfferAffiliateMQF, offerID+constants.Separator+affID, mqf)

	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Offer Affiliate")
	}
	defer pipeline.Close()

}

func SaveCounters(offerID string, affID string) {
	log.Println("From setting Counters", offerID, "OfferID", affID, "affID")
	key := offerID + constants.Separator + affID
	//create counters
	pipeline := db.RedisTranxnClient.Pipeline()
	pipeline.HSet(constants.TotalConversionCount, key, 0)
	pipeline.HSet(constants.SentConversionCount, key, 0)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving counters :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Counters")
	}
	defer pipeline.Close()

}

func SaveAdvertiserWhiteListing(advertiser string, ips []string) {
	err := db.RedisTranxnClient.SAdd(constants.AdvertiserWhiteListIPS+constants.Separator+advertiser, ips).Err()
	if err != nil {
		fmt.Println("Redis error while saving advertiser :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Advertiser")
	}

}

func SaveClick(transactionID string, offerID string, affiliateID string, AffSub string, AffSub2 string) {

	pipeline := db.RedisTranxnClient.Pipeline()
	pipeline.HSet(constants.Transactions, transactionID, offerID+constants.Separator+affiliateID+constants.Separator+AffSub+constants.Separator+AffSub2)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving click :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Click")
	}
	defer pipeline.Close()

}

func SaveClickBatch(tranxns []string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	retry := 1
	for _, tranxn := range tranxns {
		keys := strings.Split(tranxn, constants.ObjectSeparator)
		pipeline.HSet(constants.Transactions, keys[0], keys[1])
	}
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving click batch :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Click Batch. Retry:-"+strconv.Itoa(retry))
		if retry < 4 {
			SaveClickBatch(tranxns)
			retry++
		}
	}
	defer pipeline.Close()
}

func SaveUnSentPostBacks(transactionID string, offer_id string, affiliate_id string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	key := offer_id + constants.Separator + affiliate_id
	pipeline.HSet(constants.ConvertedTransactionIDHash, transactionID, constants.OfferDefault)
	pipeline.HIncrBy(constants.TotalConversionCount, key, 1)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving unsent postbacks :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Un Sent Postbacks")
	}
	defer pipeline.Close()
}

func SaveSentPostBacks(transactionID string, offer_id string, affiliate_id string) {

	pipeline := db.RedisTranxnClient.Pipeline()
	key := offer_id + constants.Separator + affiliate_id
	pipeline.HSet(constants.ConvertedTransactionIDHash, transactionID, constants.OfferDefault)
	pipeline.HSet(constants.SentTransactionIDHash, transactionID, constants.OfferDefault)
	pipeline.HIncrBy(constants.TotalConversionCount, key, 1)
	pipeline.HIncrBy(constants.SentConversionCount, key, 1)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving sent postback :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Sent Postbacks")
	}
	defer pipeline.Close()

}
func SaveUnSentPostBackBatch(tranxns []string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	for _, tranxn := range tranxns {
		keys := strings.Split(tranxn, constants.ObjectSeparator)
		pipeline.HSet(constants.ConvertedTransactionIDHash, keys[0], constants.OfferDefault)
		pipeline.HIncrBy(constants.TotalConversionCount, keys[1], 1)
	}
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving unsent postback batch :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Un Sent Postback Batch")
	}
	defer pipeline.Close()
}

func SaveSentPostBackBatch(tranxns []string) {
	pipeline := db.RedisTranxnClient.Pipeline()
	for _, tranxn := range tranxns {
		keys := strings.Split(tranxn, constants.ObjectSeparator)
		pipeline.HSet(constants.ConvertedTransactionIDHash, keys[0], constants.OfferDefault)
		pipeline.HSet(constants.SentTransactionIDHash, keys[0], constants.OfferDefault)
		pipeline.HIncrBy(constants.TotalConversionCount, keys[1], 1)
		pipeline.HIncrBy(constants.SentConversionCount, keys[1], 1)
	}
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving sent postback batch :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save Sent Postbacks Batch")
	}
	defer pipeline.Close()

}

//rotations
func SaveRotationOfferByGEOOnTranxn(geo string, offID string) {
	err := db.RedisTranxnClient.HSet(constants.RotationOfferGeoStack, geo, offID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save rotaion Offer GEO")
	}
}

func SaveRotationOfferByGroupOnTranxn(offID string, geo string, rotationOfferID string) {
	err := db.RedisTranxnClient.HSet(constants.RotationOfferGroupStack, offID+constants.Separator+geo, rotationOfferID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Save rotaion Offer Group")
	}
}

//validations
func ValidateTransactionID(transactionId string) bool {
	result := db.RedisTranxnClient.HGet(constants.Transactions, transactionId).Val()
	if len(result) == 0 {
		fmt.Println("TransactionID is not present in transaction ID Hash")
		return false
	}
	return true

}
func ValidateConvertedTransactionID(transactionId string) bool {

	result := db.RedisTranxnClient.HGet(constants.ConvertedTransactionIDHash, transactionId).Val()

	if len(result) == 0 {
		fmt.Println("TransactionID is not present in converted transaction IDHash")
		return false
	}
	return true

}
func ValidateSentTransactionID(transactionId string) bool {

	result := db.RedisTranxnClient.HGet(constants.SentTransactionIDHash, transactionId).Val()

	if len(result) == 0 {
		log.Print("Transaction Id Not Present in sent conversions", transactionId)
		return false
	}
	return true

}

func ValidateTransactionIDForPostback(transactionID string, offerID string) (bool, string, string) {

	transaction := db.RedisTranxnClient.HGet(constants.Transactions, transactionID).Val()
	transactionsSplitted := strings.Split(transaction, constants.Separator)
	affiliateID := constants.AffiliateDefault
	if len(transactionsSplitted) > 1 {
		affiliateID = transactionsSplitted[1]
	}
	offType := db.RedisTranxnClient.HGet(constants.RedisOfferType, offerID).Val()
	converted := db.RedisTranxnClient.HGet(constants.ConvertedTransactionIDHash, transactionID).Val()
	if len(converted) == 0 {
		log.Print("Not a converted transactionID:", transactionID)
		return false, offType, affiliateID
	}
	log.Print("Offer Type is returned:-", offType)
	return true, offType, affiliateID

}

func ValidateAdvertiserIP(advertiser string, ip string) bool {

	key := constants.AdvertiserWhiteListIPS + constants.Separator + advertiser
	ips := db.RedisTranxnClient.SMembers(key).Val()
	if len(ips) == 0 || strings.Compare(ips[0], "0.0.0.0") == 0 {
		return true
	} else {
		valid := false
		for _, ipeach := range ips {
			if strings.Compare(ip, ipeach) == 0 {
				valid = true
				break
			}
		}
		return valid
	}
}

func ValidateAdveriserIPWithTransaction(transactionID string, ip string) (bool, string) {
	transaction := db.RedisTranxnClient.HGet(constants.Transactions, transactionID).Val()
	transactionsSplitted := strings.Split(transaction, constants.Separator)
	offerID := constants.OfferDefault
	if len(transactionsSplitted) > 0 {
		offerID = transactionsSplitted[0]
	}
	advertiser := db.RedisTranxnClient.HGet(constants.RedisOfferAdvertiser, offerID).Val()
	valid := ValidateAdvertiserIP(advertiser, ip)
	return valid, offerID
}

func ValidateAdveriserWithOfferID(offerID string, ip string) bool {
	advertiser := db.RedisTranxnClient.HGet(constants.RedisOfferAdvertiser, offerID).Val()
	valid := ValidateAdvertiserIP(advertiser, ip)
	return valid
}

func GetOfferGroupOnTranxn(offID string) string {
	group := db.RedisTranxnClient.HGet(constants.RedisOfferGroup, offID).Val()
	return group
}
func GetOfferTypeOnTranxn(offerID string) string {
	offerType := db.RedisTranxnClient.HGet(constants.RedisOfferType, offerID).Val()
	return offerType
}
func GetOfferAffiliatePostbackTemplate(offID string, affID string) string {
	template := db.RedisTranxnClient.HGet(constants.RedisOfferAffiiatePostback, offID+constants.Separator+affID).Val()
	return template
}

func GetMQFByOfferAndAffiliate(offerID string, affiliateID string) string {
	mqf := db.RedisTranxnClient.HGet(constants.RedisOfferAffiliateMQF, offerID+constants.Separator+affiliateID).Val()
	return mqf
}

func GetMQFByAffiliate(affiliateID string) string {
	mqf := db.RedisTranxnClient.HGet(constants.RedisAffiliateMQF, affiliateID).Val()
	return mqf
}

func GetTemplateByAffiliateID(affiliateID string) string {
	template := db.RedisTranxnClient.HGet(constants.RedisAffiliateMediaTemplate, affiliateID).Val()
	return template
}

func GetConversionData(offer_id string, affiliate_id string) (string, string) {

	key := offer_id + constants.Separator + affiliate_id
	conversion_count := db.RedisTranxnClient.HGet(constants.TotalConversionCount, key).Val()
	sent_conversion_count := db.RedisTranxnClient.HGet(constants.SentConversionCount, key).Val()

	if len(sent_conversion_count) == 0 || len(conversion_count) == 0 {
		fmt.Println("Error getting data for conversion_count")
		//create counters
		SaveCounters(offer_id, affiliate_id)
		return constants.Zero, constants.Zero
	}
	return conversion_count, sent_conversion_count
}

func GetRotatedConversionData(offer_id string) (string, string) {

	key := offer_id + constants.Separator + constants.TRACKER_MEDIA
	conversion_count := db.RedisTranxnClient.HGet(constants.TotalConversionCount, key).Val()
	sent_conversion_count := db.RedisTranxnClient.HGet(constants.SentConversionCount, key).Val()
	if len(sent_conversion_count) == 0 || len(conversion_count) == 0 {
		fmt.Println("Error getting data for rotation conversion_count")
		//create mqf for offer and affiliate
		SaveOfferAffiliateMQFOnPostback(offer_id, constants.TRACKER_MEDIA, "0.3")
		//create counters
		SaveCounters(offer_id, constants.TRACKER_MEDIA)
		return constants.Zero, constants.Zero

	}
	return conversion_count, sent_conversion_count
}

func GetTransaction(transactionID string) (string, string) {
	transaction := db.RedisTranxnClient.HGet(constants.Transactions, transactionID).Val()
	transactionsSplitted := strings.Split(transaction, constants.Separator)
	offerID := constants.OfferDefault
	affiliateID := constants.AffiliateDefault
	if len(transactionsSplitted) > 1 {
		offerID = transactionsSplitted[0]
		affiliateID = transactionsSplitted[1]
	}
	return offerID, affiliateID
}

func GetTransactionWithAffiliate(transactionID string) (string, string) {
	transaction := db.RedisTranxnClient.HGet(constants.Transactions, transactionID).Val()
	transactionsSplitted := strings.Split(transaction, constants.Separator)
	affSub := ""
	affSub2 := ""
	if len(transactionsSplitted) > 2 {
		affSub = transactionsSplitted[2]
		affSub2 = transactionsSplitted[3]
	}
	return affSub, affSub2
}

func RemoveAffiliate(affID string) {
	err := db.RedisTranxnClient.SRem(constants.ExhaustedOfferStack, affID).Err()
	if err != nil {
		fmt.Print("Redis query remove affiliate error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove Affiliate")
	}
}

func RemoveMQF(offID string, affID string) {
	err := db.RedisTranxnClient.SRem(constants.ExhaustedOfferStack, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove mqf error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisTranxn", "Remove MQF")
	}
}
