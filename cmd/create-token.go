package main

import (
	"flag"
	"fmt"

	"github.com/sycomancy/glasnik/infra"
	"go.mongodb.org/mongo-driver/bson"
)

// create default OTP token in DB
func main() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
	flag.Parse()
	_, err := infra.InsertDocument("tokens", bson.D{{Key: "token", Value: "12345"}})
	if err != nil {
		fmt.Printf("Unable to create token %d", err)
		return
	}

	fmt.Println("Created new token")
}
