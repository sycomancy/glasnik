package main

import (
	"flag"
	"fmt"

	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)
	token := flag.String("token", "", "Token to be created")
	flag.Parse()
	_, err := infra.InsertDocument("tokens", bson.D{{Key: "token", Value: token}})
	if err != nil {
		fmt.Printf("Unable to create token %d", err)
		return
	}

	fmt.Println("Created new token")
}
