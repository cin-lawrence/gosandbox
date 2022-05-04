package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	BrokerURI      string `mapstructure:"BROKER_URI"`
	DatabaseURI    string `mapstructure:"DATABASE_URI"`
	Port           int    `mapstructure:"PORT"`
	CeleryTaskName string `mapstructure:"CELERY_TASK_NAME"`
}

var Config config = LoadConfig()

func setDefaults(v *viper.Viper) {
	v.SetDefault("CELERY_TASK_NAME", "fibonacci")
}

func newViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath("./configs/api")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	setDefaults(v)

	return v
}

func LoadConfig() config {
	v := newViper()

	var cfg config
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
