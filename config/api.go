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
	godotenv.Load()

	return API{
		Mongo: Mongo{
			Dsn: fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=%s",
				os.Getenv("MONGODB_USER"),
				os.Getenv("MONGODB_PASS"),
				os.Getenv("MONGODB_HOST"),
				os.Getenv("MONGODB_APPNAME"),
			),
		},
	}
}
