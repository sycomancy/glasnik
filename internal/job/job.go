package job

import (
	"fmt"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/njuskalo"
	"github.com/sycomancy/glasnik/internal/types"
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

func (j *Job) Start() error {
	itemCh := make(chan []types.AdEntry)
	go njuskalo.FetchEntry(j.Url, itemCh, j.client)

	for items := range itemCh {
		fmt.Println(items)
	}

	return nil
}

func (j *Job) persistEntries(entries []types.AdEntry) {
	models := make([]mongo.WriteModel, 0)
	for _, entry := range entries {
		model :=  
		// append(models, )
	}
	// models := []mongo.WriteModel{
	// 	mongo.NewReplaceOneModel().SetFilter(bson.D{{"name", "Cafe Tomato"}}).
	// 		SetReplacement(Restaurant{Name: "Cafe Zucchini", Cuisine: "French"}),
	// 	mongo.NewUpdateOneModel().SetFilter(bson.D{{"name", "Cafe Zucchini"}}).
	// 		SetUpdate(bson.D{{"$set", bson.D{{"name", "Zucchini Land"}}}}),
	// }
	// infra.BulkReplace("entries")
}
