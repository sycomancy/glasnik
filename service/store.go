package service

import (
	"context"
	"sync"
	"time"

	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	jobCollection            = "jobs"
	locationResultCollection = "locationResult"
)

type FetchJobEntity struct {
	Id              primitive.ObjectID `bson:"_id"`
	StartTime       time.Time          `bson:"start_time,omitempty"`
	EndTime         time.Time          `bson:"end_time,omitempty"`
	LocationsInQeue []string           `bson:"locations_queue,omitempty"`
	Completed       bool               `bson:"completed"`
}

type LocationResultEntity struct {
	Id        primitive.ObjectID `bson:"_id"`
	JobId     primitive.ObjectID `bson:"jobId"`
	LastPage  int                `bson:"lastPage"`
	Completed bool               `bson:"completed"`
	RawPages  []string           `bson:"rawPages"`
}

type Storer struct {
	locationPageMu sync.RWMutex
	jobMu          sync.RWMutex
}

func NewStorer() *Storer {
	return &Storer{}
}

func (s *Storer) CreateFetchJob(job *FetchJob) error {
	model := FetchJobEntity{
		Id:              job.Id,
		StartTime:       job.StartTime,
		EndTime:         job.EndTime,
		LocationsInQeue: job.locationsInQueue,
	}
	_, err := infra.InsertDocument(jobCollection, model)
	return err
}

func (s *Storer) GetAllRunningJobs() ([]*FetchJob, error) {
	var models []*FetchJobEntity
	infra.FindDocuments(jobCollection, bson.D{}, &models)

	jobs := make([]*FetchJob, len(models))
	for _, job := range models {
		jobs = append(jobs, &FetchJob{
			storer:           s,
			locationsInQueue: job.LocationsInQeue,
			Id:               job.Id,
			StartTime:        job.StartTime,
			EndTime:          job.EndTime,
		})
	}
	return jobs, nil
}

func (s *Storer) GetJobByID(id string) (*FetchJob, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	model := &FetchJobEntity{}
	if err != nil {
		return nil, err
	}
	coll := infra.GetCollectionByName(jobCollection)
	err = coll.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(model)

	if err != nil {
		return nil, err
	}
	return &FetchJob{
		storer:           s,
		locationsInQueue: model.LocationsInQeue,
		Id:               model.Id,
		StartTime:        model.StartTime,
		EndTime:          model.EndTime,
	}, nil
}

func (s *Storer) RemoveLocationsFromJobQueue(jobID primitive.ObjectID, locationIDS []string) *error { // Use $pull here
	s.jobMu.Lock()
	defer s.jobMu.Unlock()

	filter := bson.D{{Key: "_id", Value: jobID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "locations_queue", Value: bson.D{{Key: "$in", Value: locationIDS}}}}}}
	infra.UpdateDocument(jobCollection, filter, update)
	return nil
}

func (s *Storer) StoreResultsForLocationPage(jobID primitive.ObjectID, result *LocationPageResult, location *LocalityEntry, completed bool, lastPage int) *error {
	s.locationPageMu.Lock()
	defer s.locationPageMu.Unlock()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "location", Value: bson.D{
				{Key: "id", Value: location.Id},
				{Key: "title", Value: location.Attributes.Title},
			}},
			{Key: "jobId", Value: jobID},
			{Key: "completed", Value: completed},
			{Key: "lastPage", Value: lastPage},
		},
		},
		{Key: "$push", Value: bson.D{
			{Key: "entries", Value: bson.D{{Key: "$each", Value: result.items}}},
		}},
	}

	filter := bson.D{{Key: "location.id", Value: location.Id}}
	infra.UpsertDocument(locationResultCollection, filter, update)
	return nil
}
