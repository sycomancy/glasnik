package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "ads-fetcher",
})

var baseURL = "https://www.njuskalo.hr/prodaja-stanova?geo[locationIds]="

func main() {
	infra.NewIncognitoClient(nil)
	logg.Info("Starting ads fetcher")
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	var locations []*njuskalo.LocalityEntry
	infra.FindDocuments("locality", bson.D{}, &locations)
	worker := infra.NewWorker[*njuskalo.LocalityEntry](2)

	logg.Infof("Started work for %d locations \n", len(locations))

	worker.Produce(func(producerStream chan<- *njuskalo.LocalityEntry, stopFn func()) {
		for _, loc := range locations {
			producerStream <- loc
		}
		stopFn()
	})

	numOfFailed := 0

	worker.Consume(func(val *njuskalo.LocalityEntry) {
		fullURL := fmt.Sprintf("%s%s", baseURL, val.Id)
		client := infra.NewIncognitoClient(
			[]time.Duration{3 * time.Second, 6 * time.Second, 15 * time.Second, 30 * time.Second},
		)
		ads, err := njuskalo.FetchAds(fullURL, client)
		if err != nil {
			logg.Warnf("failed to get data for location %s %w", val.Attributes.Title, err)
			numOfFailed += 1
		}

		if len(ads) == 0 {
			logg.Warnf("check this %s", fullURL)
		}

		logg.Infof("got %d ads for %s \n", len(ads), val.Attributes.Title)
	})

	worker.Wait(func(data *njuskalo.LocalityEntry) {
		fmt.Printf("Work done. Failed jobs for %d locations", numOfFailed)
	})
}
