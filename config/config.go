package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init() {
	viper.SetConfigType("env")

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	viper.AutomaticEnv()

	log.Printf("[INFO] Config loaded %s", viper.AllSettings())
}
