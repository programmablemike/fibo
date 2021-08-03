// Implements an intermediate cache of Fibonacci values using Postgres
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/programmablemike/fibo/internal/tracing"
	log "github.com/sirupsen/logrus"
	pg "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type CacheEntry struct {
	gorm.Model
	Ordinal uint64 `gorm:"index"` // The fibonacci ordinal N
	Value   string // The fibonacci value - we use string to represent arbitrary precision
}

func (c CacheEntry) String() string {
	return fmt.Sprintf("CacheEntry<%s %s>", fibonacci.Uint64ToString(c.Ordinal), c.Value)
}

// Cache implements a PostgresDB cache for pre-computed ordinal values
type Cache struct {
	context     context.Context
	db          *gorm.DB
	initialized bool
}

// NewCache creates a new cache with persistent database connection
func NewCache(dsn string) *Cache {
	log.Debugf("Connecting to postgres with DSN=%s", dsn)
	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{
		// This turns off the default logging which is too verbose for records that don't exist
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Errorf("Failed to connect to database: %s", err)
	}
	log.Info("Successfully connected to database.")
	cache := &Cache{
		context:     context.Background(),
		db:          db,
		initialized: false,
	}
	if err := cache.init(); err != nil {
		log.Errorf("Failed to initialize the database: %s", err)
	}
	return cache
}

func (c *Cache) SetContext(ctx context.Context) {
	c.context = ctx
}

// init the cache database
func (c *Cache) init() error {
	log.Info("Initializing the database...")
	if c.initialized {
		log.Warning("Cannot re-initiliaze the database... skipping.")
		return nil
	}
	if err := c.initWaitForDatabase(); err != nil {
		log.Fatal("Timed out while waiting for database to become available")
	}
	if err := c.initTables(); err != nil {
		log.Error("Failed to initialize the database schema.")
		return err
	}
	log.Info("Database initialized successfully!")
	c.initialized = true // Mark the database as initialized
	return nil
}

func (c *Cache) initWaitForDatabase() error {
	// Wait 20 seconds for the database to come online
	timeoutAt := time.Now().Add(20 * time.Second)
	for {
		if time.Now().After(timeoutAt) {
			log.Fatal("Failed to connect to the database within 20 seconds")
			return fmt.Errorf("failed to connect to the database")
		}
		log.Info("Trying to connect to database...")
		db, err := c.db.DB()
		if err != nil {
			time.Sleep(1 * time.Second) // Wait a second to give the DB more time
			continue                    // Keep looping
		}
		err = db.Ping()
		if err != nil {
			time.Sleep(1 * time.Second) // Wait a second to give the DB more time
			continue                    // Keep looping
		}
		break // Break out of the loop
	}
	log.Info("Database connection succeeded.")
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
	if span := tracing.StartSpanFromContext(c.context, "clear"); span != nil {
		defer span.Finish()
	}
	log.Info("Clearing the database.")
	// Deletes all cache entries
	// Note that this only "tombstones" the entries in Gorm by adding a "deleted_at" timestamp
	c.db.Where("1 = 1").Delete(&CacheEntry{})
	return nil
}

func (c *Cache) Write(ordinal uint64, value *fibonacci.Number) error {
	if span := tracing.StartSpanFromContext(c.context, "write"); span != nil {
		defer span.Finish()
	}
	entry := &CacheEntry{
		Ordinal: ordinal,
		Value:   value.String(),
	}
	c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(entry)
	log.Debugf("Wrote cache entry for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	return nil
}

func (c *Cache) Read(ordinal uint64) (*fibonacci.Number, error) {
	if span := tracing.StartSpanFromContext(c.context, "read"); span != nil {
		defer span.Finish()
	}
	entry := new(CacheEntry)
	result := c.db.Where("ordinal = ?", ordinal).First(entry)
	if result.Error != nil {
		log.Debugf("Failed to retrieve cache entry for ordinal=%s: %v", fibonacci.Uint64ToString(ordinal), result.Error)
		return fibonacci.NewNumber(-1), result.Error
	}
	log.Debugf("Successfully retrieved cached value for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	v, ok := fibonacci.NewNumberFromDecimalString(entry.Value)
	if !ok {
		err := fmt.Errorf("failed to convert %s to a *fibonacci.Number", entry.Value)
		log.Error(err)
		return fibonacci.NewNumber(-1), err
	}
	log.Debugf("Read cache entry for ordinal=%s", fibonacci.Uint64ToString(ordinal))
	return v, nil
}
