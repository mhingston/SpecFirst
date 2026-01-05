package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ProjectName string            `mapstructure:"project_name"`
	Protocol    string            `mapstructure:"protocol"`
	Language    string            `mapstructure:"language"`
	Framework   string            `mapstructure:"framework"`
	CustomVars  map[string]string `mapstructure:"custom_vars"`
	Constraints map[string]string `mapstructure:"constraints"`
}

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
