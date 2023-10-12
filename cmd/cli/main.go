package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/job"
)

var url = "https://www.njuskalo.hr/prodaja-stanova?geo%5BlocationIds%5D=2698"

func main() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
	flag.Parse()
	job, err := job.NewJob(url)
	if err != nil {
		fmt.Println("unable to start job", err)
		os.Exit(-1)
	}

	job.Start()
	// err, job := job.NewFetchJob("", 1)

	// for _, url := range urls {
	// 	result, err := njuskalo.Fetch(url, client)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Printf("Got %d results for url: %s \n", len(result), url)
	// 		for _, r := range result {
	// 			fmt.Println(r.Title, r.Id, r.Price)
	// 		}
	// 	}
	// }
}
