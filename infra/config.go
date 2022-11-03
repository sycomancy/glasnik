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
			DB_URL: os.Getenv("DB_URL"),
		}
	}
}

type config struct {
	DB_URL string
}

var Config config
