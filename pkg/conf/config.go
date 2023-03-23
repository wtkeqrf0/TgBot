package conf

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	Prod *bool `yaml:"prod" env-required:"true"`
	Bot  struct {
		Secret     string `yaml:"secret" env-required:"true"`
		WeatherKey string `yaml:"weather_key" env-required:"true"`
	} `yaml:"bot"`
	Listen struct {
		Host string `yaml:"host" env-default:"127.0.0.1"`
		Port string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Postgres struct {
		Username string `yaml:"username" env-required:"true"`
		DBName   string `yaml:"db_name" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		Host     string `yaml:"host" env-default:"127.0.0.1"`
		Port     string `yaml:"port" env-default:"5432"`
	} `yaml:"postgres"`
}

var inst *Config
var once sync.Once

// GetConfig builds the configuration file in golang type and returns it
func GetConfig() *Config {
	once.Do(func() {
		inst = &Config{}

		if err := cleanenv.ReadConfig("configs/config.yml", inst); err != nil {
			help, _ := cleanenv.GetDescription(inst, nil)
			logrus.Info(help)
			logrus.Fatalf("error occurred while reading config file: %v", err)
		}
		if *inst.Prod {
			inst.Postgres.Host = "postgres"
		}
	})

	return inst
}
