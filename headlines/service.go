package headlines

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	mongoSession *mgo.Session
	httpClient   *http.Client
	listURL      string
	conceptURL   string
}

const (
	DATABASE   = "upp-store"
	COLLECTION = "content"
)

func NewHeadlineService(connStr string, listURL string, conceptURL string) service {

	session, err := mgo.Dial(connStr)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	return service{
		mongoSession: session,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		listURL:    listURL,
		conceptURL: conceptURL,
	}
}

func (s *service) getHeadlinesByUUID(UUIDs []string) ([]headlineOutput, error) {
	c := s.mongoSession.DB(DATABASE).C(COLLECTION)

	result := headlineOutput{}
	output := make([]headlineOutput, len(UUIDs), len(UUIDs))
	outputOrder := map[string]int{}

	for i, e := range UUIDs {
		outputOrder[e] = i
	}

	iter := c.Find(bson.M{"uuid": bson.M{
		"$in": UUIDs,
	}}).Select(bson.M{"_id": 0, "uuid": 1, "title": 1, "standfirst": 1}).Iter()

	for iter.Next(&result) {
		index := outputOrder[result.UUID]
		output[index] = result
	}
	err := iter.Close()

	return output, err
}

func (s *service) getHeadlinesByList(listUUID string) ([]headlineOutput, error) {
	resp, err := s.httpClient.Get(s.listURL + listUUID)
	if err != nil {
		return nil, err
	}

	list := List{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&list)
	if err != nil {
		return nil, err
	}

	var UUIDs []string
	for _, e := range list.Items {
		UUIDs = append(UUIDs, e.ID[strings.LastIndex(e.ID, "/")+1:])
	}

	logrus.Debugf("UUIDs found for list %s: %v", listUUID, UUIDs)

	return s.getHeadlinesByUUID(UUIDs)
}

func (s *service) getHeadlinesByConcept(conceptUUID string) ([]headlineOutput, error) {
	resp, err := s.httpClient.Get(s.conceptURL + conceptUUID)
	if err != nil {
		return nil, err
	}

	var items []ListItem
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&items)
	if err != nil {
		return nil, err
	}

	var UUIDs []string
	for _, e := range items {
		UUIDs = append(UUIDs, e.ID[strings.LastIndex(e.ID, "/")+1:])
	}

	logrus.Debugf("UUIDs found for concept %s: %v", conceptUUID, UUIDs)

	return s.getHeadlinesByUUID(UUIDs)
}

func (s *service) getFlashBriefingByUUID(UUIDs []string) ([]FlashBriefingItem, error) {
	c := s.mongoSession.DB(DATABASE).C(COLLECTION)

	result := []FlashBriefingItem{}

	err := c.Find(bson.M{"uuid": bson.M{
		"$in": UUIDs,
	}}).Select(bson.M{"_id": 0, "uuid": 1, "publishedDate": 1, "title": 1, "standfirst": 1}).All(&result)

	return result, err
}

func (s *service) getFlashBriefingForList(listUUID string) ([]FlashBriefingItem, error) {
	resp, err := s.httpClient.Get(s.listURL + listUUID)
	if err != nil {
		return nil, err
	}

	list := List{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&list)
	if err != nil {
		return nil, err
	}

	var UUIDs []string
	for _, e := range list.Items {
		UUIDs = append(UUIDs, e.ID[strings.LastIndex(e.ID, "/")+1:])
	}

	logrus.Debugf("UUIDs found for list %s: %v", listUUID, UUIDs)

	return s.getFlashBriefingByUUID(UUIDs)
}
