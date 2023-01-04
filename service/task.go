package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	locationWorkersNum = 3
	baseURL            = "https://www.njuskalo.hr/prodaja-stanova?geo[locationIds]="
)

var flogg = logrus.WithFields(logrus.Fields{
	"ctx": "fetchSchedule",
})

type FetchJob struct {
	storer           *Storer
	Id               primitive.ObjectID
	StartTime        time.Time
	EndTime          time.Time
	locationsInQueue []string
}

func NewFetchTask() *FetchJob {
	return &FetchJob{
		storer:    NewStorer(),
		Id:        primitive.NewObjectID(),
		StartTime: time.Now(),
	}
}

func (j *FetchJob) Run() error {
	flogg.Info("started schedule")

	locations := GetAllLocalityEntries()
	locationIDS := []string{}
	for _, l := range locations {
		locationIDS = append(locationIDS, l.Id)
	}

	j.locationsInQueue = locationIDS
	j.storer.CreateFetchJob(j)

	worker := infra.NewWorker[*LocalityEntry](locationWorkersNum)
	worker.Produce(func(producerStream chan<- *LocalityEntry, stopFn func()) {
		for _, loc := range locations {
			producerStream <- loc
		}
		stopFn()
	})

	worker.Consume(func(l *LocalityEntry) {
		// s.fetchAdsForLocation(l)
		// ads, err := njuskalo.FetchAds(fullURL, client)
		// if err != nil {
		// 	logg.Warnf("failed to get data for location %s %w", val.Attributes.Title, err)
		// }

		// if len(ads) == 0 {
		// 	logg.Warnf("check this %s", fullURL)
		// }
	})
	worker.Wait(func(data *LocalityEntry) {
		flogg.Info("done!!!")
	})
	return nil
}

func (s *FetchJob) fetchAdsForLocation(loc *LocalityEntry) {
	_ = fmt.Sprintf("%s%s", baseURL, loc.Id)
	client := infra.NewIncognitoClient([]time.Duration{3 * time.Second, 6 * time.Second, 15 * time.Second, 30 * time.Second})
	flogg.Infof("processing location %s \n", loc.Attributes.Title)
	location := NewLocation(client)
	a, err := location.GetPageHTML("", 0, client)
	if err != nil {
		flogg.Error(err)
	}
	flogg.Info(a)
	// items := make([]types.AdEntry, 0)
	// hasMorePage := true
	// page := 1

	// for hasMorePage {
	// 	pageItems, err := njuskalo
	// }
}

// ################# DB Models #####################

// func (s *FetchJob) upsertModelInDB() {
// 	filter := bson.D{{Key: "_id", Value: s.Id}}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "start_time", Value: s.CreateTime},
// 			{Key: "end_time", Value: s.EndTime},
// 		},
// 		},
// 	}
// 	infra.UpsertDocument(scheduleCollection, filter, update)
