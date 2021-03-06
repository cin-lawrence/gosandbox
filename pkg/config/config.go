package config

import (
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type config struct {
	Dev  bool `mapstructure:"DEV"`
	Test bool `mapstructure:"TEST"`

	BrokerURI   string `mapstructure:"BROKER_URI"`
	DatabaseURI string `mapstructure:"DATABASE_URI"`
	Port        int    `mapstructure:"PORT"`

	AccessTokenExpirationTime  int64  `mapstructure:"ACCESS_TOKEN_EXPIRATION_TIME"`
	RefreshTokenExpirationTime int64  `mapstructure:"REFRESH_TOKEN_EXPIRATION_TIME"`
	AccessTokenSecret          string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret         string `mapstructure:"SECRET_TOKEN_SECRET"`

	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	CeleryTaskName string `mapstructure:"CELERY_TASK_NAME"`
}

var Config config = LoadConfig()

func setDefaults(v *viper.Viper) {
	v.SetDefault("DEV", false)
	v.SetDefault("TEST", false)
	v.SetDefault("ACCESS_TOKEN_EXPIRATION_TIME", 15*60)
	v.SetDefault("REFRESH_TOKEN_EXPIRATION_TIME", 7*24*60*60)
	v.SetDefault("ACCESS_TOKEN_SECRET", "paragona")
	v.SetDefault("REFRESH_TOKEN_SECRET", "paragonf")
	v.SetDefault("CELERY_TASK_NAME", "fibonacci")
}

func newViper() *viper.Viper {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can't get caller information")
	}
	v := viper.New()
	v.AddConfigPath("./configs/api")
	v.AddConfigPath(filepath.Join(path, "../../../configs/test"))
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
