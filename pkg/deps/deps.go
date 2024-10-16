package deps

import (
	"github.com/vit0rr/short-spot/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Deps struct {
	Config   config.Config
	DBClient *mongo.Client
}

func New(config config.Config, mgClient *mongo.Client) *Deps {
	return &Deps{
		Config:   config,
		DBClient: mgClient,
	}
}
