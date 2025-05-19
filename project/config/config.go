package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`

	AppPort    string `mapstructure:"APP_PORT"`
	AppBaseURL string `mapstructure:"APP_BASE_URL"`

	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306") // Default MySQL port
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_BASE_URL", "http://localhost:8080")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			log.Println("INFO: Config file .env not found. Using environment variables and defaults.")
		} else {

			log.Printf("ERROR: Failed to read config file: %v\n", err)
			return Config{}, err // Return empty config and the error
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("ERROR: Unable to unmarshal config: %v\n", err)
		return Config{}, err
	}

	if config.WeatherAPIKey == "" {
		log.Println("WARNING: WEATHER_API_KEY is not set in the configuration.")

	}

	log.Println("Configuration loaded.")
	return config, nil
}
