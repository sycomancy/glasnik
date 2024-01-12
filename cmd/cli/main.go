package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/job"
)

var urls = []string{"https://www.njuskalo.hr/prodaja-stanova/novi-zagreb-zapad"}

func main() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
	flag.Parse()

	job, err := job.NewDetailsJob()
	if err != nil {
		fmt.Println("unable to start job", err)
		os.Exit(-1)
	}

	job.FetchDetails()

	// for _, url := range urls {
	// 	job, err := job.NewJob(url)
	// 	fmt.Println(job)
	// 	if err != nil {
	// 		fmt.Println("unable to start job", err)
	// 		os.Exit(-1)
	// 	}

	// 	fmt.Println("started job for ", url)
	// 	job.FetchEntries()
	// }
}
