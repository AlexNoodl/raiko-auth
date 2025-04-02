package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github/alexnoodl/raiko-auth/pkg/logger"
	"os"
)

type Config struct {
	Port     string
	MongoURI string
	JWTKey   string
	DBName   string
	Logger   *logrus.Logger
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	logger := logger.SetupLogger()

	if err != nil {
		logger.Warn("Error loading .env file: ", err)
	}

	cfg := &Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		JWTKey:   getEnv("JWT_KEY", "secret"),
		DBName:   getEnv("DB_NAME", "auth"),
		Logger:   logger,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
