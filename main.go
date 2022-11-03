package main

import (
	"flag"
	"fmt"

	"github.com/sycomancy/glasnik/client"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

func main() {
	// load config
	infra.LoadConfig()
	// Connect to mongodb
	infra.MongoConnect(infra.Config.DB_URL)

	listenAddr := flag.String("listenAddr", ":3000", "Listen address the server is running")
	flag.Parse()
	svc := NewLoggingService(NewTokenValidatorService(&adsFetcher{}))

	// Just for testing
	go func() {
		client := client.NewClient(":3333", "/result-webhook", "http://localhost:3000/api/request-njuska")

		go func() {
			client.Run()
		}()

		_, err := client.SendRequest(&types.Request{
			Filter:      "https://www.njuskalo.hr/prodaja-stanova?geo%5BlocationIds%5D=2698",
			Token:       "121345",
			CallbackURL: "http://localhost:3333/result-webhook",
		})

		if err != nil {
			fmt.Println(err)
		}

		for {
			data := <-client.Data
			fmt.Println("got results from service", data.RequestID)
		}

	}()

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()
}
