package internal

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	PostgresURL string
	Port        string
	JWTSecret   []byte
	TestMode    bool
}

func LoadAppConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return &AppConfig{
		PostgresURL: os.Getenv("GOOSE_DBSTRING"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   []byte(os.Getenv("JWT_SECRET")),
		TestMode:    os.Getenv("API_TEST") == "true",
	}
}
