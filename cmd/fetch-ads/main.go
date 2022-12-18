package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
	"go.mongodb.org/mongo-driver/bson"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "ads-fetcher",
})

func main() {
	logg.Info("Starting ads fetcher")
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	locations := make([]*njuskalo.LocalityEntry, 0)
	infra.FindDocuments("locality", bson.D{}, locations)
	fmt.Println(locations)
}

// var filter = "https://www.njuskalo.hr/prodaja-stanova/zagreb"

// type aggregator struct {
// 	client *infra.IncognitoClient
// }

// func (a *aggregator) FetchAndPersist(filters []string) {
// 	results, err := njuskalo.FetchAds(filter, a.client)
// 	if err != nil {
// 		logg.Error(err)
// 	}

// 	a.persist(results)
// 	logg.Info("got: ", len(results), " results for filter: ", filter)
// }

// func (a *aggregator) persist(results []types.AdEntry) {
// 	entries := make([]interface{}, len(results))
// 	for i := range results {
// 		entries[i] = results[i]
// 	}
// 	_ = infra.InsertDocuments("results", entries)
// }

// func Run() {
// 	aggregator := &aggregator{client: infra.NewIncognitoClient(nil)}
// 	aggregator.FetchAndPersist([]string{filter})
// }
