package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	APIKey              string        `mapstructure:"API_KEY"`
	DeliveryAPI_URL     string        `mapstructure:"DELIVERY_API_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	// Load configuration from the specified path
	// This is a placeholder for actual implementation, e.g., using viper or another library
	// For now, we will return a default config for demonstration purposes
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
