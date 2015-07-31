package matchmaking

import "gopkg.in/fly/config.v1"

// Config describes the config object..
type Config struct {
	TeamPlayerCount int      `yaml:"team_player_count"`
	TeamCount       int      `yaml:"teams"`
	Plugins         []string `yaml:"plugins"`
}

func NewConfigFromNamespace(organization string, system string) (c *Config, err error) {
	var cfgNS = config.Namespace{
		Organization: "fly",
		system:       "matchmaking",
	}
	err = cfgNS.Load(c)
	if err != nil {
		return
	}
	return
}

// NewConfig reads configuration from environment and returns a config object.
func NewConfig() (c *Config, err error) {
	c, err = NewConfigFromNamespace("fly", "matchmaking")
	return
}
