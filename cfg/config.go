package cfg

import "github.com/fly/config"

// Config describes the config object..
type Config struct {
	TeamSize     int          `yaml:"team_size"`
	TeamCount    int          `yaml:"teams"`
	Checks       []string     `yaml:"checks"`
	IRCBotConfig IRCBotConfig `yaml:"irc_bot"`
}

// IRCBotConfig holds the config options of the matchmaking IRC bot
type IRCBotConfig struct {
	Server string `yaml:"server"`
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
	// password to connect to the IRC network
	Password string `yaml:"password"`
	Channels []struct {
		// name of the channel to join
		Channel string `yaml:"channel"`
		// password for entering the channel
		Password string `yaml:"password"`
	} `yaml:"channels"`
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
