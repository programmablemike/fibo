package cache

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

type Cache struct {
	Store *pgx.Conn
}

// NewCache returns a new cache instance
func NewCache() *Cache {
	conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost:5432/fibo")
	if err != nil {
		log.Fatalf("Failed to connect to the cache: %s", err)
	}
	defer conn.Close(context.Background())
	return &Cache{Store: conn}
}

// Init creates the database and tables if they don't exist
func (c *Cache) Init() error {
	if c.Store.IsClosed() {
		return fmt.Errorf("The postgres connection is closed")
	}

	return nil
}
