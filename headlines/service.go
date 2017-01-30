package headlines

import (
	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	mongoSession *mgo.Session
}

const (
	DATABASE   = "upp-store"
	COLLECTION = "content"
)

func newHeadlineService(connStr string) service {
	session, err := mgo.Dial(connStr)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	return service{
		mongoSession: session,
	}
}

func (s *service) getHeadlines(UUIDs []string) []headlineOutput {
	log.Info("getHeadlines")
	c := s.mongoSession.DB("upp-store").C("content")

	result := []headlineOutput{}

	err := c.Find(bson.M{"uuid": bson.M{
		"$in": UUIDs,
	}}).Select(bson.M{"_id": 0, "uuid": 1, "title": 1, "standfirst": 1}).All(&result)

	if err != nil {
		panic(err)
	}
	log.Info(result)

	return result
}
