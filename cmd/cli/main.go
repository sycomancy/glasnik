package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/job"
)

var urls = []string{"https://www.njuskalo.hr/prodaja-stanova/novi-zagreb-istok"}

func main() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
	flag.Parse()
	for _, url := range urls {
		job, err := job.NewJob(url)
		if err != nil {
			fmt.Println("unable to start job", err)
			os.Exit(-1)
		}

		fmt.Println("started job for ", url)
		job.FetchEntries()
	}
}
