package main

import (
	"github.com/sycomancy/glasnik/aggregator"
	"github.com/sycomancy/glasnik/infra"
)

func main() {
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	aggregator.Run()
}
