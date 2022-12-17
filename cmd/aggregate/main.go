package main

import (
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
)

func main() {
	infra.LoadConfig()
	infra.MongoConnect(infra.Config.DB_URL)

	metaBuilder := njuskalo.NewNjuskaloMeta(infra.NewIncognitoClient(nil))

	err := metaBuilder.RebuildLocalityMeta([]string{"1264", "1263", "1262", "1261", "1260", "1259", "1256", "1247", "1248", "1249", "1250", "1251", "1252", "1253", "1254", "1255", "1257", "1258"})
	if err != nil {
		panic(err)
	}
}
