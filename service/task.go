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

	locationService := NewLocationService()

	worker := infra.NewWorker[*LocalityEntry](locationWorkersNum)
	worker.Produce(func(producerStream chan<- *LocalityEntry, stopFn func()) {
		for _, loc := range locations {
			producerStream <- loc
		}
		stopFn()
	})

	worker.Consume(func(l *LocalityEntry) {
		j.fetchAdsForLocation(l, locationService)
	})

	worker.Wait(func(data *LocalityEntry) {
		flogg.Info("done!!!")
	})
	return nil
}

func (j *FetchJob) fetchAdsForLocation(loc *LocalityEntry, service *LocationService) {
	flogg.Infof("processing location %s %s \n", loc.Id, loc.Attributes.Title)

	client := infra.NewIncognitoClient([]time.Duration{3 * time.Second, 6 * time.Second, 15 * time.Second, 30 * time.Second})

	locationPageResult := make(chan *LocationPageResult, 1000)
	service.GetLocationPages(loc, locationPageResult, client)
	result := <-locationPageResult
	fmt.Println("Got result for", loc.Id, loc.Attributes.Title, result.completed, result.page)

	j.storer.StoreResultsForLocationPage(j.Id, result, loc, result.completed, result.page)
	j.storer.RemoveLocationsFromJobQueue(j.Id, []string{loc.Id})
}
