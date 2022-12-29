package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO(sycomancy): implement DB interface, remove Mongo dependency

var (
	scheduleCollection = "tasks"
	locationWorkersNum = 3
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
	return nil
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

func (s *FetchTask) usertModelInDB() {
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
