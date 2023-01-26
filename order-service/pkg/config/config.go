package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	DBUrl         string `mapstructure:"DB_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InDockerComposeEnv() *Config {
	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		return c
	}

	db_user := os.Getenv("DB_USER")
	db_port := os.Getenv("DB_PORT")
	db_pswd := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	c.DBUrl = db_user + ":" + db_pswd + "@tcp(" + db_host + ":" + db_port + ")/" + db_name
	c.ProductSvcUrl = "product-service:50052"
	return c
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

	err = viper.Unmarshal(c)
	if err != nil {
		return c
	}

	return c
}
