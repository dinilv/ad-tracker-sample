package v1

import (
	"context"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	helpers "github.com/adcamie/adserver/helpers"
	api "github.com/micro/micro/api/proto"
)

type Impression struct {
	ImpressionListener
}

type ImpressionListener interface {
	Track2(context.Context, *api.Request, *api.Response) error
}

func (tracker *Impression) Tracker(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to track impression:")

	//process cookies
	requestParams := make(map[string]string)
	//prcocess request and parameters
	helpers.ParseGetRequest(requestParams, req)
	//add escape path
	requestParams[constants.URL] = "/aff_i?" + requestParams[constants.URL]
	helpers.ParseGetHeader(requestParams, req)

	requestParams["click_url"] = requestParams[constants.URL]
	requestParams[constants.URL] = ""

	//set cookie if its empty
	if len(requestParams[constants.ADSAUCE_ID]) == 0 {
		userid := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		values := []string{constants.ADSAUCE_ID + "=" + userid + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;"}
		//set cookies in response
		rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
		//save to message for further processing
		requestParams[constants.ADSAUCE_ID] = userid
	}

	//put details to subscriber for async processing
	requestParams["activity"] = "1"
	helpers.Subscribe(constants.ImpressionTopic, requestParams)

	//return sucesss for impression always
	rsp.StatusCode = 200
	return nil
}
