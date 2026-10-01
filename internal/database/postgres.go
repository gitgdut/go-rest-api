package database

import (
	"context"

	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Unableto parse DATABASE_URL: %v", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Printf("Successfully connected to PostgreSQL database")
	return pool, nil
}
