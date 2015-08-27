package cfg

import "github.com/fly/config"

// Config describes the config object..
type Config struct {
	TeamSize  int          `yaml:"team_size"`
	TeamCount int          `yaml:"teams"`
	Checks    []string     `yaml:"checks"`
	IRC       IRCBotConfig `yaml:"irc_bot"`
}

// IRCBotConfig holds the config options of the matchmaking IRC bot
type IRCBotConfig struct {
	// hostname:port of IRC server to connect to.
	Server string `yaml:"server"`
	// Name of the bot
	Name   string `yaml:"name"`
	UseTLS bool   `yaml:"use_tls"`
	// password to connect to the IRC network
	Password      string `yaml:"password"`
	CommandPrefix string `yaml:"command_prefix"`
	Channels      []struct {
		// name of the channel to join
		Name string `yaml:"Name"`
		// password for entering the channel
		Password string `yaml:"password"`
	} `yaml:"channels"`
}

// LoadConfig loads configuration from file and returns the object
func LoadConfig() (c *Config, err error) {
	var cfgNS = config.Namespace{
		Organization: "fly",
		System:       "matchmaking",
	}

	var c = &Config{}
	err = cfgNS.Load(&c)
	if err != nil {
		return
	}

	return
}
