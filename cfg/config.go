package cfg

// Config describes the config object..
type Config struct {
	TeamSize  int      `yaml:"team_size"`
	TeamCount int      `yaml:"teams"`
	Checks    []string `yaml:"checks"`
}
