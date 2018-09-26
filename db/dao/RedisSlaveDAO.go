package dao

/* This class provides functions that are used at time of click for validating received click & postback against different criteria
Also it provides functions to support selection of rotating to offer.
*/
import (
	"fmt"
	"strings"

	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
)

//Rotation criteria for an offer eligibility
func ValidateExhaustedOffer(offerID string) bool {
	result := db.RedisSlaveClient.SIsMember(constants.ExhaustedOfferStack, offerID).Val()
	if result {
		fmt.Println("Exhausted Offer.", offerID, result)
		return true
	}
	return false

}

func ValidateExhaustedOfferAffiliate(offerID string, affiliateID string) bool {
	result := db.RedisSlaveClient.SIsMember(constants.ExhaustedOfferAffiliateStack, offerID+constants.Separator+affiliateID).Val()
	if result {
		fmt.Println("Exhausted Offer and afiliate.", offerID, affiliateID, result)
		return true
	}
	return false

}

//Rotation criteria for an offer eligibility
func ValidateExhaustedOfferHash(offerID string) bool {
	result := db.RedisSlaveClient.HGet(constants.ExhaustedOfferHash, offerID).Val()
	if len(result) == 0 || strings.Compare(result, constants.Zero) != 0 {
		fmt.Println("Exhausted Offer.", offerID, result)
		return false
	} else {
		return true
	}
}

func ValidateExhaustedOfferAffiliateHash(offerID string, affiliateID string) bool {
	result := db.RedisSlaveClient.HGet(constants.ExhaustedOfferAffiliateHash, offerID+constants.Separator+affiliateID).Val()
	if len(result) == 0 || strings.Compare(result, constants.Zero) != 0 {
		fmt.Println("Exhausted Offer and Affiliate.", offerID, result)
		return false
	} else {
		return true
	}
}
func ValidateOfferCountry(offerID string, countryID string) bool {

	key := constants.OfferGeoTargetting + constants.Separator + offerID
	countries := db.RedisSlaveClient.SMembers(key).Val()
	fmt.Println("Offer countries", countries)
	if len(countries) == 0 || len(countryID) == 0 || strings.Compare(countries[0], "ALL") == 0 {
		fmt.Println("Valid Country")
		return true
	} else {
		valid := false
		for _, country := range countries {
			fmt.Println(country, countryID)
			if strings.Compare(country, countryID) == 0 || strings.Compare(country, "ALL") == 0 {
				valid = true
				fmt.Println("Valid country")
				break

			}
		}
		fmt.Println("Valid", valid)
		return valid
	}

}

func ValidateBlackListIP(offerID string, ip string) bool {
	return true
}

func ValidateOperator(offerID string, carrier string) bool {
	return true
}

//query required data of offer for further processing
func GetTemplateByOfferID(offerID string) string {
	template := db.RedisSlaveClient.HGet(constants.RedisOfferAdvertiserTemplate, offerID).Val()
	return template
}
func GetOperatorTemplateByOfferID(offerID string, operator string) string {
	template := db.RedisSlaveClient.HGet(constants.RedisOfferOperatorTemplate, offerID+constants.ObjectSeparator+operator).Val()
	return template
}

func GetOfferGroup(offID string) string {
	group := db.RedisSlaveClient.HGet(constants.RedisOfferGroup, offID).Val()
	return group
}

func GetOfferType(offerID string) string {
	offerType := db.RedisSlaveClient.HGet(constants.RedisOfferType, offerID).Val()
	return offerType
}

//query for rotations
func GetRotationOfferByGEO(countryCode string) string {
	offer := db.RedisSlaveClient.HGet(constants.RotationOfferGeoStack, countryCode).Val()
	return offer
}
func GetRotationMOOfferByGEO(countryCode string) string {
	offer := db.RedisSlaveClient.HGet(constants.RotationMOOfferGeoStack, countryCode).Val()
	return offer
}
func GetRotationOfferByOfferID(offerID string, countryCode string) string {
	offer := db.RedisSlaveClient.HGet(constants.RotationOfferGroupStack, offerID+constants.Separator+countryCode).Val()
	return offer
}
