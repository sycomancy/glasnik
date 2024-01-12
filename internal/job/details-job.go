package job

import (
	"context"
	"fmt"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/njuskalo"
	"github.com/sycomancy/glasnik/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DetailsJob struct {
	client *infra.IncognitoClient
}

func NewDetailsJob() (*DetailsJob, error) {
	job := &DetailsJob{}
	job.client = infra.NewIncognitoClient(nil)
	return job, nil
}

func (j *DetailsJob) FetchDetails() error {
	var result []types.AdInDb
	coll := infra.DBClient.Database("colusa").Collection("entries")
	collDetails := infra.DBClient.Database("colusa").Collection("detailedEntries")

	filter := bson.D{{}}

	cursor, err := coll.Find(nil, filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	entriesCh := make(chan types.FetchEntryDetailsResult)
	go njuskalo.FetchEntryDetails(result, entriesCh, j.client)
	for entry := range entriesCh {
		j.persistEntry(entry, collDetails)
	}
	return nil
}

func (j *DetailsJob) persistEntry(item types.FetchEntryDetailsResult, coll *mongo.Collection) {
	fmt.Println("persisting entry", item.Ad.ID)
	filter := bson.D{{Key: "_id", Value: item.Ad.ID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "ad", Value: item.Ad},
			{Key: "details", Value: item.Details},
		}},
	}

	updateOptions := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.TODO(), filter, update, updateOptions)

	if err != nil {
		panic(err)
	}

}
