// pkg/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpireHours int    `mapstructure:"JWT_EXPIRE_HOURS"`

	EmailMock bool `mapstructure:"EMAIL_MOCK"`
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("no .env file found, using environment variables: %v", err)
	}

	// Viper's Unmarshal does not include values discovered only through
	// AutomaticEnv. Read each setting directly so container environment
	// variables and local .env values behave identically.
	cfg := &Config{
		AppEnv:         viper.GetString("APP_ENV"),
		AppPort:        viper.GetString("APP_PORT"),
		DBHost:         viper.GetString("DB_HOST"),
		DBPort:         viper.GetString("DB_PORT"),
		DBName:         viper.GetString("DB_NAME"),
		DBUser:         viper.GetString("DB_USER"),
		DBPassword:     viper.GetString("DB_PASSWORD"),
		RedisAddr:      viper.GetString("REDIS_ADDR"),
		RedisPassword:  viper.GetString("REDIS_PASSWORD"),
		JWTSecret:      viper.GetString("JWT_SECRET"),
		JWTExpireHours: viper.GetInt("JWT_EXPIRE_HOURS"),
		EmailMock:      viper.GetBool("EMAIL_MOCK"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.JWTExpireHours == 0 {
		cfg.JWTExpireHours = 24
	}

	return cfg
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable TimeZone=Asia/Shanghai"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "" || c.AppEnv == "development"
}
