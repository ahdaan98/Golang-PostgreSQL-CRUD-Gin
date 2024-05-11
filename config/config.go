package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	User     string
	DBname   string
	Password string
	Port     string
	DSN      string
}

func LoadEnv() Config {
	file := ".env"
	err := godotenv.Load(file)

	if err != nil {
		log.Fatal("FAILED TO LOAD" + file)
	}

	config := Config{
		Host:     os.Getenv("host"),
		User:     os.Getenv("user"),
		DBname:   os.Getenv("dbname"),
		Password: os.Getenv("dbpassword"),
		Port:     os.Getenv("dbport"),
	}

	config.DSN = DSN(config)
	return config
}

func DSN(cfg Config) string {
	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.DBname, cfg.Password, cfg.Port)
}
