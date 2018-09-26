package dao

/*
As name suggests this class provides functions that is exclusively for reading instance critical data for Click.
Contains Writing & Removal of data to different keys.
*/
import (
	"fmt"

	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	logger "github.com/adcamie/adserver/logger"
)

//insertions
func SaveOffer(offerID string, offerType string, template string, group string, advertiser string, countryIds []string) {
	pipeline := db.RedisMasterClient.Pipeline()
	pipeline.HSet(constants.RedisOfferAdvertiserTemplate, offerID, template)
	pipeline.HSet(constants.RedisOfferType, offerID, offerType)
	pipeline.Del(constants.OfferGeoTargetting + constants.Separator + offerID)
	for _, countryID := range countryIds {
		pipeline.SAdd(constants.OfferGeoTargetting+constants.Separator+offerID, countryID)
	}
	pipeline.HSet(constants.RedisOfferGroup, offerID, group)
	//for redis-tranxn on advertiser verification
	SaveOfferOnTranxn(offerID, offerType, template, group, advertiser, countryIds)
	_, err := pipeline.Exec()
	if err != nil {
		fmt.Println("Redis error while saving offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save Offer")
	}
	defer pipeline.Close()
}

func SaveExhaustedOffer(offID string) {
	SaveExhaustedOfferOnTranxn(offID)
	err := db.RedisMasterClient.SAdd(constants.ExhaustedOfferStack, offID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save exhausted Offer")
	}
}

func SaveExhaustedOfferAffiliate(offerID string, affiliateID string) {
	SaveExhaustedOfferAffiliateOnTranxn(offerID, affiliateID)
	err := db.RedisMasterClient.SAdd(constants.ExhaustedOfferAffiliateStack, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save Exhausted Offer Affiliate")
	}
}

func SaveExhaustedOfferMap(offID string) {
	SaveExhaustedOfferMapOnTranxn(offID)
	err := db.RedisMasterClient.HSet(constants.ExhaustedOfferHash, offID, constants.Zero).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save exhausted Offer on HashMap")
	}
}

func SaveExhaustedOfferAffiliateMap(offerID string, affiliateID string) {
	SaveExhaustedOfferAffiliateMapOnTranxn(offerID, affiliateID)
	err := db.RedisMasterClient.HSet(constants.ExhaustedOfferAffiliateHash, offerID+constants.Separator+affiliateID, constants.Zero).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save Exhausted Offer Affiliate on HashMap")
	}
}

//rotation related
func SaveRotationOfferByGEO(geo string, offID string) {
	SaveRotationOfferByGEOOnTranxn(geo, offID)
	err := db.RedisMasterClient.HSet(constants.RotationOfferGeoStack, geo, offID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save rotaion Offer geo")
	}
}

func SaveRotationOfferByGroup(offID string, geo string, rotationOfferID string) {
	SaveRotationOfferByGroupOnTranxn(offID, geo, rotationOfferID)
	err := db.RedisMasterClient.HSet(constants.RotationOfferGroupStack, offID+constants.Separator+geo, rotationOfferID).Err()
	if err != nil {
		fmt.Println("Redis error while saving exhausted offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Save exhausted Offer")
	}
}

//query
func GetOfferTypeOnMaster(offerID string) string {
	offerType := db.RedisMasterClient.HGet(constants.RedisOfferType, offerID).Val()
	return offerType
}

//remove keys
func RemoveOfferInExhausted(offID string) {
	RemoveOfferInExhaustedOnTranxn(offID)
	err := db.RedisMasterClient.SRem(constants.ExhaustedOfferStack, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Remove Offer In Exhausted")
	}
}

func RemoveOfferAffiliateInExhausted(offerID string, affiliateID string) {
	RemoveOfferAffiliateInExhaustedOnTranxn(offerID, affiliateID)
	err := db.RedisMasterClient.SRem(constants.ExhaustedOfferAffiliateStack, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer affiliate in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Remove Offer Affiliate MQF")
	}
}

//remove keys
func RemoveOfferInExhaustedHash(offID string) {
	RemoveOfferInExhaustedOnTranxnHash(offID)
	err := db.RedisMasterClient.HDel(constants.ExhaustedOfferHash, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Remove Offer In Exhausted Hash")
	}
}

func RemoveOfferAffiliateInExhaustedHash(offerID string, affiliateID string) {
	RemoveOfferAffiliateInExhaustedOnTranxnHash(offerID, affiliateID)
	err := db.RedisMasterClient.HDel(constants.ExhaustedOfferAffiliateHash, offerID+constants.Separator+affiliateID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer affiliate in exhausted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Remove Offer Affiliate MQF in Exhausted Hash")
	}
}
func RemoveOffer(offID string) {
	err := db.RedisMasterClient.SRem(constants.RedisOfferAdvertiser, offID).Err()
	err = db.RedisMasterClient.SRem(constants.RedisOfferAdvertiserTemplate, offID).Err()
	err = db.RedisMasterClient.SRem(constants.RedisOfferAffiiatePostback, offID).Err()
	err = db.RedisMasterClient.SRem(constants.RedisOfferAffiliateMQF, offID).Err()
	err = db.RedisMasterClient.SRem(constants.RedisOfferGroup, offID).Err()
	err = db.RedisMasterClient.SRem(constants.RedisOfferGroup, offID).Err()
	if err != nil {
		fmt.Print("Redis query remove offer error :", err.Error())
		go logger.ErrorLogger(err.Error(), "RedisMaster", "Remove Offer")
	}
}
