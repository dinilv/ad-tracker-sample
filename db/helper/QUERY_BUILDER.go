package helper

import (
	"strconv"
)

func BigQueryBuilder(todayDate string, hour int, collection string) string {

	var query string

	switch collection {

	/*case "ClickLog":
		query = `SELECT OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID,
                 COUNT(DISTINCT MessageID) as messageid,
				 COUNT(DISTINCT(CASE WHEN Activity=2 OR Activity=13 THEN messageid END)) as Clicks,
				 COUNT(DISTINCT(CASE WHEN Activity=13 THEN messageid END)) as RotatedClicksFwd
				 FROM adcamie_tracker.click_log
				 WHERE date(Date)="` + todayDate + `" AND Hour =` + strconv.Itoa(hour) +
			     ` AND OfferID IS NOT NULL AND AffiliateID IS NOT NULL
			     GROUP BY OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID;`*/

	case "PostBackLog":
		query = `SELECT OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID,
			     COUNT(DISTINCT TransactionID) as transactionid,
				 COUNT(DISTINCT(CASE WHEN Activity=3 OR Activity=11 OR Activity=4 OR Activity=12 OR Activity=7 OR Activity=8 THEN transactionid END)) as Conversions,
                 COUNT(DISTINCT(CASE WHEN Activity=3 OR Activity=11  THEN transactionid END)) as SentConversions,
                 COUNT(DISTINCT(CASE WHEN Activity=4 OR Activity=12  THEN transactionid END)) as UnSentConversions,
                 COUNT(DISTINCT(CASE WHEN Activity=7 THEN transactionid END)) as SentRotatedConversionsFwd,
                 COUNT(DISTINCT(CASE WHEN Activity=8 THEN transactionid END)) as UnSentRotatedConversionsFwd
				 FROM adcamie_tracker.postback_log
				 WHERE date(Date)="` + todayDate + `" AND Hour =` + strconv.Itoa(hour) +
			     ` AND OfferID IS NOT NULL AND AffiliateID IS NOT NULL AND TransactionID IS NOT NULL
			     GROUP BY OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID;`

	case "PostBackLogWithoutTransaction":
		query = `SELECT OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID,
			     COUNT(*) as transactionid,
				 COUNT(DISTINCT(CASE WHEN Activity=3 OR Activity=11 OR Activity=4 OR Activity=12 OR Activity=7 OR Activity=8 THEN transactionid END)) as Conversions,
                 COUNT(DISTINCT(CASE WHEN Activity=3 OR Activity=11  THEN transactionid END)) as SentConversions,
                 COUNT(DISTINCT(CASE WHEN Activity=4 OR Activity=12  THEN transactionid END)) as UnSentConversions,
                 COUNT(DISTINCT(CASE WHEN Activity=7 THEN transactionid END)) as SentRotatedConversionsFwd,
                 COUNT(DISTINCT(CASE WHEN Activity=8 THEN transactionid END)) as UnSentRotatedConversionsFwd
				 FROM adcamie_tracker.postback_log
				 WHERE date(Date)="` + todayDate + `" AND Hour =` + strconv.Itoa(hour) +
			     ` AND OfferID IS NOT NULL AND AffiliateID IS NOT NULL AND TransactionID IS NULL
			     GROUP BY OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID;`

	case "PostEventLog":
		query = `SELECT OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID,
			     COUNT(DISTINCT MessageID) as messageid,
			     COUNT(DISTINCT(CASE WHEN Activity=5  THEN messageid END)) as sent_postevents,
                 COUNT(DISTINCT(CASE WHEN Activity=6  THEN messageid END)) as unsent_postevents,
                 COUNT(DISTINCT(CASE WHEN Activity=9  THEN messageid END)) as sent_rotated_postevents,
                 COUNT(DISTINCT(CASE WHEN Activity=10  THEN messageid END)) as unsent_rotated_postevents,
                 COUNT(DISTINCT(CASE WHEN Activity=5 OR Activity=6 OR Activity=9 OR Activity=10 THEN messageid END)) as PostEvents
				 FROM adcamie_tracker.postevent_log
				 WHERE date(Date)="` + todayDate + `" AND Hour =` + strconv.Itoa(hour) +
			     ` AND OfferID IS NOT NULL AND AffiliateID IS NOT NULL
			     GROUP BY OfferID,AffiliateID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID;`

	case "ClickCookieIDLog":
		query = `SELECT OfferID,CookieID, COUNT (CookieID) as CookieCount
				 FROM adcamie_tracker.click_log
				 WHERE OfferID IS NOT NULL AND CookieID IS NOT NULL
                 GROUP BY CookieID,OfferID
				 HAVING COUNT (CookieID) > 4
				 ORDER BY OfferID DESC;`

	case "PostbackCookieIDLog":
		query = `SELECT OfferID,CookieID, COUNT (CookieID) as CookieCount
				 FROM adcamie_tracker.postback_log
				 WHERE OfferID IS NOT NULL AND CookieID IS NOT NULL
                 GROUP BY CookieID,OfferID
				 HAVING COUNT (CookieID) > 1
				 ORDER BY OfferID DESC;`

	case "ClickLog":
		query = `SELECT OfferID,AffiliateID,hour,Date,ReceivedOfferID,ReceivedAffiliateID,
                 SUM(total_clicks) as Clicks,
                 SUM(unique_clicks) as UniqueClicks,
                 SUM(duplicate_clicks) as DuplicateClicks,
                 SUM(messageid) as messageid,
                 SUM(rotated_clicks_fwd) as RotatedClicksFwd
                 FROM (SELECT OfferID,AffiliateID,CookieID,Date,hour,ReceivedOfferID,ReceivedAffiliateID,
                       COUNT (DISTINCT MessageID) as messageid,
                       COUNT(DISTINCT(CASE WHEN Activity=13 THEN messageid END)) as rotated_clicks_fwd,
                       COUNT (CookieID) as total_clicks,
                       COUNT (DISTINCT CookieID) as unique_clicks,
                       COUNT (CookieID)- COUNT (DISTINCT CookieID)  as duplicate_clicks
                       FROM adcamie_tracker.click_log
                       WHERE date(Date)="` + todayDate + `" AND Hour =` + strconv.Itoa(hour) +
			           ` AND OfferID IS NOT NULL AND AffiliateID IS NOT NULL
			           GROUP BY OfferID,AffiliateID,CookieID,Date,Hour,ReceivedOfferID,ReceivedAffiliateID)
                 GROUP BY OfferID,AffiliateID,Date,hour,ReceivedOfferID,ReceivedAffiliateID;`

	}

	return query
}
