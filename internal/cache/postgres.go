// Implements an intermediate cache of Fibonacci values using Postgres
package cache

import (
	"fmt"

	"github.com/programmablemike/fibo/internal/fibonacci"
	log "github.com/sirupsen/logrus"
	pg "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

type CacheEntry struct {
	gorm.Model
	Ordinal uint64 `gorm:"uniqueIndex"` // The fibonacci ordinal N
	Value   string // The fibonacci value - we use string to represent arbitrary precision
}

func (c CacheEntry) String() string {
	return fmt.Sprintf("CacheEntry<%s %s>", fibonacci.Uint64ToString(c.Ordinal), c.Value)
}

// Cache implements a PostgresDB cache for pre-computed ordinal values
type Cache struct {
	db          *gorm.DB
	initialized bool
}

// NewCache creates a new cache with persistent database connection
func NewCache(dsn string) *Cache {
	log.Debugf("Connecting to postgres with DSN=%s", dsn)
	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("Failed to connect to database: %s", err)
	}
	log.Info("Successfully connected to database.")
	cache := &Cache{
		db:          db,
		initialized: false,
	}
	if err := cache.init(); err != nil {
		log.Errorf("Failed to initialize the database: %s", err)
	}
	return cache
}

// init the cache database
func (c *Cache) init() error {
	log.Info("Initializing the database...")
	if c.initialized {
		log.Warning("Cannot re-initiliaze the database... skipping.")
		return nil
	}
	if err := c.initTables(); err != nil {
		log.Error("Failed to initialize the database schema.")
		return err
	}
	log.Info("Database initialized successfully!")
	c.initialized = true // Mark the database as initialized
	return nil
}

// initSchema creates the table schema
func (c *Cache) initTables() error {
	c.db.AutoMigrate(&CacheEntry{})
	log.Info("Successfully initialized the table schemas.")
	return nil
}

func (c *Cache) Close() error {
	log.Info("Closing the database connection.")
	db, _ := c.db.DB()
	return db.Close()
}

func (c *Cache) Clear() error {
	return nil
}

func (c *Cache) Write(ordinal uint64, value *fibonacci.Number) error {
	entry := &CacheEntry{
		Ordinal: ordinal,
		Value:   value.String(),
	}
	c.db.Create(entry)
	log.Infof("Wrote cache entry for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	return nil
}

func (c *Cache) Read(ordinal uint64) (*fibonacci.Number, error) {
	entry := new(CacheEntry)
	result := c.db.Where("ordinal = ?", ordinal).First(entry)
	if result.Error != nil {
		log.Warningf("Failed to retrieve cache entry for ordinal=%s: %v", fibonacci.Uint64ToString(ordinal), result.Error)
		return fibonacci.NewNumber(-1), result.Error
	}
	log.Debugf("Successfully retrieved cached value for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	v, ok := fibonacci.NewNumberFromDecimalString(entry.Value)
	if !ok {
		err := fmt.Errorf("failed to convert %s to a *fibonacci.Number", entry.Value)
		log.Error(err)
		return fibonacci.NewNumber(-1), err
	}
	log.Infof("Read cache entry for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	return v, nil
}
