package matchmaking

import "gopkg.in/fly/config.v1"

// Config describes the config object..
type Config struct {
	TeamPlayerCount int      `yaml:"team_player_count"`
	TeamCount       int      `yaml:"teams"`
	Plugins         []string `yaml:"plugins"`
}

// NewConfig reads configuration from environment and returns a config object.
func NewConfig(organization string, system string) (c *Config, err error) {
	var cfgNS = config.Namespace{
		Organization: organization,
		system:       system,
	}
	err = cfgNS.Load(c)
	if err != nil {
		return
	}
	return
}
