package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Elasticsearch ElasticsearchConfig
	MinIO         MinIOConfig
	JWT           JWTConfig
	CORS          CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type JWTConfig struct {
	Secret string
	Expiry string
}

type CORSConfig struct {
	Origins []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "3000"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "suipic"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "suipic"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Elasticsearch: ElasticsearchConfig{
			Addresses: strings.Split(getEnv("ES_ADDRESSES", "http://localhost:9200"), ","),
			Username:  getEnv("ES_USERNAME", ""),
			Password:  getEnv("ES_PASSWORD", ""),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:    getBoolEnv("MINIO_USE_SSL", false),
			Bucket:    getEnv("MINIO_BUCKET", "suipic"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		CORS: CORSConfig{
			Origins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3001"), ","),
		},
	}

	return config, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}
