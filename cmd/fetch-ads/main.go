package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
	"go.mongodb.org/mongo-driver/bson"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "ads-fetcher",
})

func main() {
	logg.Info("Starting ads fetcher")
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	var locations []*njuskalo.LocalityEntry
	infra.FindDocuments("locality", bson.D{{"_id", "2656"}}, &locations)
	for _, loc := range locations {
		fmt.Println(loc.Attributes.Title)
	}
}

// func main() {
// 	worker := NewWorker(10)

// worker.Produce(func(producerStream chan<- int, stopFn func()) {
// 			for i := 0; i < 5; i++ {
// 					producerStream <- i
// 			}

// 			stopFn()
// 	})

// worker.Consume(func(val int) {
// 	fmt.Printf("Consumer is processing %d\n", val)
// })

// 	worker.Wait(nil)
// }
