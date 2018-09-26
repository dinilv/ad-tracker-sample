package v1

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	helper "github.com/adcamie/adserver/helpers"
	api "github.com/micro/micro/api/proto"
)

type Cv struct {
	ContentviewListener
}

type ContentviewListener interface {
	Track(context.Context, *api.Request, *api.Response) error
}

func (tracker *Cv) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to CMP Content View:")

	//cookies & request
	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()

	//req & params
	helper.ParseGetRequest(requestParams, req)

	//escape path & host & headers
	requestParams[constants.HOST] = constants.CMP_HOST
	requestParams[constants.URL] = constants.AFF_CV + requestParams[constants.URL]
	helper.ParseGetHeader(requestParams, req)
	requestParams[constants.PING_URL] = requestParams[constants.URL]
	requestParams[constants.URL] = ""

	//validate subscription
	requestParams[constants.ACTIVITY] = "63"

	//async processing
	helper.Subscribe(constants.ContentViewTopic, requestParams)

	//sucesss
	rsp.StatusCode = 200
	return nil
}
