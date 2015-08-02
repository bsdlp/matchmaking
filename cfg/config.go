package cfg

// Config describes the config object..
type Config struct {
	TeamPlayerCount int      `yaml:"team_player_count"`
	TeamCount       int      `yaml:"teams"`
	Checks          []string `yaml:"checks"`
}
