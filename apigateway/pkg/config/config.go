package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	AuthSvcUrl    string `mapstructure:"AUTH_SVC_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
	OrderSvcUrl   string `mapstructure:"ORDER_SVC_URL"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InLocalConfig() *Config {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return c
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c
	}

	return c
}

func (c *Config) InDockerComposeEnv() *Config {
	aurl := os.Getenv("DB_HOST")
	if aurl == "" {
		return c
	}

	c.AuthSvcUrl = "auth-service:50051"
	c.ProductSvcUrl = "product-service:50052"
	c.OrderSvcUrl = "order-service:50053"

	return c
}
