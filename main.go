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

	listenAddr := flag.String("listenAddr", ":3000", "Listen address the server is running")
	flag.Parse()
	svc := NewLoggingService(NewTokenValidatorService(&adsFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()
}
