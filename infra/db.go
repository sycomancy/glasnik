package infra

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// TODO(dhudek): read from config
var dbName = "colusa"

func MongoConnect(uri string) error {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"db uri": uri,
		}).Panic("failed connecting to db")
		return err
	} else {
		logrus.WithFields(logrus.Fields{
			"db uri": uri,
		}).Info("connected to db")
		return nil
	}
}

func FindDocument[T any](collection string, filter bson.D, result *T) *T {
	coll := client.Database(dbName).Collection(collection)
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		panic(err)
	}

	return result
}

func FindDocuments[T any](collection string, filter bson.D, result []*T) []*T {
	coll := client.Database(dbName).Collection(collection)
	curr, err := coll.Find(context.TODO(), filter)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		panic(err)
	}
	if err = curr.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	return result
}

func CountDocuments(collection string, filter bson.D) int64 {
	coll := client.Database(dbName).Collection(collection)
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return count
}

func InsertDocument(collection string, document any) (*mongo.InsertOneResult, error) {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertDocuments(collection string, documents []interface{}) *mongo.InsertManyResult {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateDocument(collection string, filter bson.D, update bson.D) *mongo.UpdateResult {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateDocuments(collection string, filter bson.D, update bson.D) *mongo.UpdateResult {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result
}

func DeleteDocument(collection string, filter bson.D) *mongo.DeleteResult {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return result
}

func DeleteDocuments(collection string, filter bson.D) *mongo.DeleteResult {
	coll := client.Database(dbName).Collection(collection)
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return result
}
