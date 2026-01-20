package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/razedwell/go-hand/internal/platform/logger"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	JWTAccessSecret  string
	JWTRefreshSecret string
	Timezone         string
	RedisAddr        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
}

func LoadConfig() *Config {
	err := godotenv.Load() // Loads .env file
	if err != nil {
		logger.Log.Println("No .env file found, using system environment variables")
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "user"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBName:           getEnv("DB_NAME", "backend_db"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "default_access_secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "default_refresh_secret"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          redisDB,
		Timezone:         getEnv("TIMEZONE", "UTC"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
