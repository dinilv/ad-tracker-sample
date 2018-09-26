package v1

import api "github.com/micro/micro/api/proto"

func ParseGetRequestOnAPI(req *api.Request, reqG map[string]string) {
	for _, get := range req.Get {
		for _, val := range get.Values {
			reqG[get.Key] = val
		}
	}
}
