package database

import (
	"context"
	"fmt"

	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func ConnectDB( /*ctx context.Context cfg Config*/ ) (*pgxpool.Pool, error) {
	dbConfig, err := loadConfig()
	if err != nil {
		// log.Fatalf("Failed to load config: %v", err)
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config, err := pgxpool.ParseConfig(dbConfig.DSN)
	if err != nil {
		// log.Fatalf("Failed to parse database configuration: %v", err)
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	config.MaxConns = dbConfig.MaxConns
	config.MinConns = dbConfig.MinConns
	config.MaxConnLifetime = dbConfig.MaxConnLifetime
	config.MaxConnIdleTime = dbConfig.MaxConnIdleTime

	// pool, err := pgxpool.NewWithConfig(ctx, config)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Use a context with timeout for Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close() // should i have this here???
		return nil, fmt.Errorf("failed to test the connection: %w", err)
	}
	return pool, nil
}

func loadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed godotenv.Load")
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	dbURL := ConcatDSN()

	return Config{
		DSN:             dbURL,
		MaxConns:        10,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}, nil
}

func ConcatDSN() string {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)
}
