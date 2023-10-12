package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/job"
)

var url = "https://www.njuskalo.hr/prodaja-stanova/novi-zagreb-zapad"

func main() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
	flag.Parse()
	job, err := job.NewJob(url)
	if err != nil {
		fmt.Println("unable to start job", err)
		os.Exit(-1)
	}

	job.FetchEntries()
}
