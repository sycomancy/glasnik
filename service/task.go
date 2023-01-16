package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	locationWorkersNum = 5
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

	// TODO(sycomancy): change this to for result := range locationPageResult
	for {
		result := <-locationPageResult
		if result.Err != nil {
			flogg.Fatal("-----", result.Err)
			return
		}

		flogg = flogg.WithFields(logrus.Fields{"completed": result.Completed, "result_count": len(result.Items), "error": result.Err, "job_id": j.Id.Hex(), "loc_id": loc.Id, "loc": loc.Title})
		flogg.Info("fetch adds for location")

		j.storer.StoreResultsForLocationPage(j.Id, result, loc, result.Completed, result.Page)
		if result.Completed {
			j.storer.RemoveLocationsFromJobQueue(j.Id, []string{loc.Id})
			return
		}
	}
}
