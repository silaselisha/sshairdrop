package util

import (
	"github.com/spf13/viper"
)

type config struct {
	DbName              string `mapstructure:"DB_NAME"`
	DbUri               string `mapstructure:"DB_URI"`
	SenderEmailAddress  string `mapstructure:"SENDER_EMAIL_ADDRESS"`
	SenderEmailName     string `mapstructure:"SENDER_EMAIL_NAME"`
	SenderEmailPassword string `mapstructure:"SENDER_EMAIL_PASSWORD"`
	TokenSecretKey      string `mapstructure:"TOKEN_SECRET_KEY"`
	ServerAddress       string `mapstructure:"SERVER_ADDRESS"`
}

func Load(path string) (config *config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return
}
