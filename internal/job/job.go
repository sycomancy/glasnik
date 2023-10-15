package job

import (
	"fmt"
	"time"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/njuskalo"
	"github.com/sycomancy/glasnik/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	Url    string
	client *infra.IncognitoClient
	jobId  interface{}
}

type FetchResult struct {
	types.AdEntry
	JobID interface{}
}

func NewJob(url string) (*Job, error) {
	job := &Job{Url: url}
	job.client = infra.NewIncognitoClient(nil)
	return job, nil
}

func (j *Job) FetchEntries() error {
	j.createJobEntry()
	entriesCh := make(chan []types.AdEntry)
	go njuskalo.FetchEntries(j.Url, entriesCh, j.client)
	for entries := range entriesCh {
		j.persistEntries(entries)
	}
	return nil
}

func (j *Job) createJobEntry() {
	date := time.Now()
	inserted, err := infra.InsertDocument("jobs", bson.D{
		{
			Key:   "started",
			Value: date.Unix(),
		},
	})
	if err != nil {
		panic("unable to insert job in db")
	}
	j.jobId = inserted.InsertedID
}

func (j *Job) persistEntries(entries []types.AdEntry) {
	models := make([]mongo.WriteModel, 0)
	for _, entry := range entries {
		e := bson.D{
			{
				Key:   "jobId",
				Value: j.jobId,
			},
			{
				Key:   "add",
				Value: entry,
			},
		}
		filter := bson.D{{Key: "slug", Value: entry.Id}, {Key: "jobId", Value: j.jobId}}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(bson.M{"$set": e}).SetUpsert(true)
		models = append(models, model)
	}
	r, err := infra.BulkWrite("entries", models)
	if err != nil {
		fmt.Println("unable to insert data to DB", err)
		return
	}

	fmt.Println("upserted: ", r.UpsertedCount, " modified: ", r.ModifiedCount)
}
