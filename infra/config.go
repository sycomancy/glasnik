package infra

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
		panic("Failed to load environment variables.")
	} else {
		Config = config{
			DB_URL:                os.Getenv("DB_URL"),
			NJUSKALO_BEARER_TOKEN: os.Getenv("NJUSKA_BEARER_TOKEN"),
		}
	}
}

type config struct {
	DB_URL                string
	NJUSKALO_BEARER_TOKEN string
}

var Config config
