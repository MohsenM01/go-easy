package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	AppName         string `mapstructure:"APP_NAME"`
	Environment     string `mapstructure:"ENV"`
	Port            string `mapstructure:"PORT"`
	DbConnection    string `mapstructure:"DB_CONNECTION"`
	RedisConnection string `mapstructure:"REDIS_CONNECTION"`
	RabbitMQ        struct {
		URL string `mapstructure:"RABBITMQ_URL"`
	} `mapstructure:",squash"`
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		viper.SetDefault("PORT", 8080)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Not Found!")
		}

		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Can not read the file : %v", err)
		}
	})

	return config
}
