package config

import "github.com/hashicorp/hcl/v2/hclsimple"

// Config is the top-level config
type Config struct {
	Server Server `hcl:"server,block"`
	API    API    `hcl:"api,block"`
}

// Server related config
type Server struct {
	BindAddr   string `hcl:"bind_addr,attr"`
	LogLevel   string `hcl:"log_level,attr"`
	CtxTimeout int    `hcl:"ctx_timeout,attr"`
}

// GetConfig returns a config from an hcl file
func GetConfig(path string) (Config, error) {
	config := Config{}
	err := hclsimple.DecodeFile(path, nil, &config)
	return config, err
}

// DefaultConfig returns a default config
func DefaultConfig() Config {
	return Config{
		Server: Server{
			BindAddr:   ":8080",
			LogLevel:   "INFO",
			CtxTimeout: 5,
		},
		API: GetDefaultAPIConfig(),
	}
}
