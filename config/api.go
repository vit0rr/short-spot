package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// API related config
type API struct {
	Mongo Mongo `hcl:"mongo,block"`
}

type Mongo struct {
	Dsn string `hcl:"dsn,attr"`
}

func GetDefaultAPIConfig() API {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return API{
		Mongo: Mongo{
			Dsn: fmt.Sprintf("mongodb://%s:%s@localhost:27017",
				os.Getenv("MONGODB_USER"),
				os.Getenv("MONGODB_PASS"),
			),
		},
	}
}
