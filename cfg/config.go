package cfg

import "github.com/fly/config"

// Config describes the config object..
type Config struct {
	TeamSize  int      `yaml:"team_size"`
	TeamCount int      `yaml:"teams"`
	Checks    []string `yaml:"checks"`
}

// LoadConfig loads configuration from file and returns the object
func LoadConfig() (c *Config, err error) {
	var cfgNS = config.Namespace{
		Organization: "fly",
		Systemn:      "matchmaking",
	}

	var c = &Config{}
	err = cfgNS.Load(&c)
	if err != nil {
		return
	}

	return
}
