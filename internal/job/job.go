package job

import (
	"fmt"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/njuskalo"
	"github.com/sycomancy/glasnik/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	Url    string
	client *infra.IncognitoClient
}

func NewJob(url string) (*Job, error) {
	job := &Job{Url: url}
	job.client = infra.NewIncognitoClient(nil)
	return job, nil
}

func (j *Job) FetchEntries() error {
	entriesCh := make(chan []types.AdEntry)
	go njuskalo.FetchEntries(j.Url, entriesCh, j.client)
	for entries := range entriesCh {
		j.persistEntries(entries)
	}

	return nil
}

func (j *Job) persistEntries(entries []types.AdEntry) {
	models := make([]mongo.WriteModel, 0)
	for _, entry := range entries {
		model := mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "id", Value: entry.Id}}).SetUpdate(bson.M{"$set": entry}).SetUpsert(true)
		models = append(models, model)
	}
	r, err := infra.BulkWrite("entries", models)
	if err != nil {
		fmt.Println("unable to insert data to DB", err)
		return
	}

	fmt.Println("upserted: ", r.UpsertedCount, " modified: ", r.ModifiedCount)
}
