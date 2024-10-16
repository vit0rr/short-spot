package config

// API related config
type API struct {
	Postgres Postgres `hcl:"postgres,block"`
}

type Postgres struct {
	Dsn string `hcl:"dsn,attr"`
}

func GetDefaultAPIConfig() API {
	return API{
		Postgres: Postgres{
			Dsn: "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable",
		},
	}
}
