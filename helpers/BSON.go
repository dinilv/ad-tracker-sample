package v1

import "gopkg.in/mgo.v2/bson"

//Generic Function To Convert String Array to Bson Data [Input type needs to be array of string]

func ConvertToBson(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}