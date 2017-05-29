package mongo

import (
	"github.com/simlenghong/yelp-fusion/yelp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)

func CreateRespSearchBusinessMin(pRespSearchBusinessMin  *yelp.RespSearchBusinessMin) (bool, error) {
	lSession, lErr := mgo.Dial("localhost:27017")
	if lErr != nil {
		return false, lErr
	}
	defer lSession.Close()

	lCollection := lSession.DB("yelp").C("RespSearchBusinessMin")
	lErr = lCollection.Insert(pRespSearchBusinessMin)
	if lErr != nil {
		return false, lErr
	}
	return true, lErr
}

func GetRespSearchBusinessMin(pTerm string, pLocation string) (*yelp.RespSearchBusinessMin, error) {
	var lRespSearchBusinessMin = new(yelp.RespSearchBusinessMin)

	lSession, lErr := mgo.Dial("localhost:27017")
	if lErr != nil {
		return lRespSearchBusinessMin, lErr
	}
	defer lSession.Close()

	lCollection := lSession.DB("yelp").C("RespSearchBusinessMin")
	lErr = lCollection.Find(bson.M{"term": pTerm, "location": pLocation}).One(&lRespSearchBusinessMin)
	if lErr != nil {
		return lRespSearchBusinessMin, lErr
	}
	return lRespSearchBusinessMin, lErr
}
