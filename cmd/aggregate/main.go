package main

import (
	"fmt"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
)

func main() {
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	metaBuilder := njuskalo.NewNjuskaloMeta(infra.NewIncognitoClient(nil))
	data, err := metaBuilder.RebuildLocalityMeta([]int{1153})
	if err != nil {
		panic(err)
	}
	fmt.Println(data.Data)
}
