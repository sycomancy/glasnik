package main

import (
	"flag"

	"github.com/sycomancy/glasnik/infra"
)

func main() {
	// load config
	infra.LoadConfig()
	// Connect to mongodb
	infra.MongoConnect(infra.Config.DB_URL)
	flag.Parse()
}
