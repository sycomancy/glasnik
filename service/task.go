package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	locationWorkersNum = 1
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
	client := infra.NewIncognitoClient([]time.Duration{3 * time.Second, 6 * time.Second, 15 * time.Second, 30 * time.Second})

	locationPageResult := make(chan *LocationPageResult)
	go service.GetLocationPages(loc, locationPageResult, client)
	result := <-locationPageResult

	if result.err != nil {
		flogg.Fatal("-----", result.err)
		return
	}

	j.storer.StoreResultsForLocationPage(j.Id, result, loc, result.completed, result.page)

	if result.completed {
		j.storer.RemoveLocationsFromJobQueue(j.Id, []string{loc.Id})
	}

	flogg = flogg.WithFields(logrus.Fields{"completed": result.completed, "result_count": len(result.items), "error": result.err, "job_id": j.Id.Hex(), "loc_id": loc.Id, "loc": loc.Attributes.Title})
	flogg.Info("fetch adds for location")
}
