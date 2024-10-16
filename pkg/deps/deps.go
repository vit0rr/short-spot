package deps

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vit0rr/short-spot/config"
)

type Deps struct {
	Config config.Config
	DBPool *pgxpool.Pool
}

func New(config config.Config, dbPool *pgxpool.Pool) *Deps {
	return &Deps{
		Config: config,
		DBPool: dbPool,
	}
}
