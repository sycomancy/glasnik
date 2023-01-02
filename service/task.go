package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO(sycomancy): implement DB interface, remove Mongo dependency

var (
	scheduleCollection = "tasks"
	locationWorkersNum = 3
	baseURL            = "https://www.njuskalo.hr/prodaja-stanova?geo[locationIds]="
)

var flogg = logrus.WithFields(logrus.Fields{
	"ctx": "fetchSchedule",
})

type FetchTaskModel struct {
	Id              primitive.ObjectID `bson:"_id"`
	StartTime       time.Time          `bson:"start_time,omitempty"`
	EndTime         time.Time          `bson:"end_time,omitempty"`
	LocationsInQeue []string           `bson:"locations_qeue,omitempty"`
	Completed       bool               `bson:"completed"`
}

type FetchTask struct {
	Id               primitive.ObjectID
	CreateTime       time.Time
	EndTime          time.Time
	locationsInQueue []string
}

// TODO(sycomancy): implement NewFetchTaskFromID
func NewFetchTask() *FetchTask {
	s := &FetchTask{
		Id:         primitive.NewObjectID(),
		CreateTime: time.Now(),
	}
	s.createModel()
	return s
}

func (s *FetchTask) Run() error {
	flogg.Info("started schedule")

	locations := GetAllLocalityEntries()
	lIDS := []string{}
	for _, l := range locations {
		lIDS = append(lIDS, l.Id)
	}
	s.locationsInQueue = lIDS

	worker := infra.NewWorker[*LocalityEntry](locationWorkersNum)

	worker.Produce(func(producerStream chan<- *LocalityEntry, stopFn func()) {
		for _, loc := range locations {
			producerStream <- loc
		}
		stopFn()
	})

	worker.Consume(func(loc *LocalityEntry) {
		s.fetchAdsForLocation(loc)
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

func (s *FetchTask) fetchAdsForLocation(loc *LocalityEntry) {
	_ = fmt.Sprintf("%s%s", baseURL, loc.Id)
	_ = infra.NewIncognitoClient([]time.Duration{3 * time.Second, 6 * time.Second, 15 * time.Second, 30 * time.Second})
	flogg.Infof("processing location %s \n", loc.Attributes.Title)

	items := make([]types.AdEntry, 0)
	hasMorePage := true
	page := 1

	for hasMorePage {
		pageItems, err := njuskalo
	}
}

// ################# DB Models #####################

func (s *FetchTask) createModel() {
	infra.InsertDocument(scheduleCollection, FetchTaskModel{
		Id:              s.Id,
		StartTime:       s.CreateTime,
		EndTime:         s.EndTime,
		LocationsInQeue: s.locationsInQueue,
	})
}

func (s *FetchTask) upsertModelInDB() {
	filter := bson.D{{Key: "_id", Value: s.Id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "start_time", Value: s.CreateTime},
			{Key: "end_time", Value: s.EndTime},
		},
		},
	}
	infra.UpsertDocument(scheduleCollection, filter, update)
}
