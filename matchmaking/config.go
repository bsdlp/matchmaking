package matchmaking

import "github.com/fly/config"

// Config describes the config object..
type Config struct {
	TeamPlayerCount int      `yaml:"team_player_count"`
	TeamCount       int      `yaml:"teams"`
	Checks          []string `yaml:"checks"`
}

// NewConfig reads configuration from environment and returns a config object.
func NewConfig() (c *config.Config, err error) {
	c, err = config.NewConfigFromNamespace("fly", "matchmaking")
	return
}
