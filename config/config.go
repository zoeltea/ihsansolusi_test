package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	AppPort    string `envconfig:"APP_PORT" default:"8080"`
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName     string `envconfig:"DB_NAME" default:"accounts_db"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	LogLevel   string `envconfig:"LOG_LEVEL" default:"info"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load(configPath)

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error processing env config: %w", err)
	}

	return &cfg, nil
}

func NewDatabaseConnection(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Verify the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
