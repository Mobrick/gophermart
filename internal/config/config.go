package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr              string
	FlagLogLevel             string
	FlagDBConnectionAddress  string
	FlagAccrualSystemAddress string
}

func MakeConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.FlagRunAddr, "a", ":8080", "address to run server")
	flag.StringVar(&config.FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&config.FlagDBConnectionAddress, "d", "", "database connection address")
	flag.StringVar(&config.FlagAccrualSystemAddress, "r", "", "points calculation system address")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.FlagLogLevel = envLogLevel
	}

	if envDBConnectionAddress := os.Getenv("DATABASE_DSN"); envDBConnectionAddress != "" {
		config.FlagDBConnectionAddress = envDBConnectionAddress
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.FlagAccrualSystemAddress = envAccrualSystemAddress
	}
	return config
}
