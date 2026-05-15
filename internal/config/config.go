package config

import (
	"errors"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ConfigFile string
	Cookie     string
	Wordlist   string
	Threads    int
	Salt       string
	Digest     string
	Verbose    bool
	Timeout    time.Duration
	Output     string
	JSON       bool
	Quiet      bool
}

func Defaults(v *viper.Viper) {
	v.SetDefault("threads", 50)
	v.SetDefault("salt", "cookie-session")
	v.SetDefault("digest", "sha1")
	v.SetDefault("timeout", 0)
}

func Setup(v *viper.Viper, configFile string) error {
	v.SetEnvPrefix("BFLASK")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/bflask")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if configFile == "" && errors.As(err, &notFound) {
			return nil
		}
		return err
	}

	return nil
}

func Load(v *viper.Viper) Config {
	return Config{
		ConfigFile: v.ConfigFileUsed(),
		Cookie:     v.GetString("cookie"),
		Wordlist:   v.GetString("wordlist"),
		Threads:    v.GetInt("threads"),
		Salt:       v.GetString("salt"),
		Digest:     v.GetString("digest"),
		Verbose:    v.GetBool("verbose"),
		Timeout:    v.GetDuration("timeout"),
		Output:     v.GetString("output"),
		JSON:       v.GetBool("json"),
		Quiet:      v.GetBool("quiet"),
	}
}
