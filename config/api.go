package config

// API related config
type API struct {
	Mongo Mongo `hcl:"mongo,block"`
}

type Mongo struct {
	Dsn string `hcl:"dsn,attr"`
}

func GetDefaultAPIConfig() API {
	return API{
		Mongo: Mongo{
			Dsn: "mongodb://docker:docker@localhost:27017",
		},
	}
}
