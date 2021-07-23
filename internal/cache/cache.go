package cache

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pg "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

type CacheEntry struct {
	gorm.Model
	Ordinal int `gorm:"uniqueIndex"` // The fibonacci ordinal N
	Result  int // The fibonacci value
}

func (c CacheEntry) String() string {
	return fmt.Sprintf("CacheEntry<%d %d>", c.Ordinal, c.Result)
}

// Cache implements a PostgresDB cache for pre-computed ordinal values
type Cache struct {
	db          *gorm.DB
	initialized bool
}

// NewCache creates a new cache with persistent database connection
func NewCache(user string, password string, addr string, database string) *Cache {
	dsn := "postgres://fibo:averysecurepasswordshouldgohere@localhost:15432/fibo"
	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %e", err)
	}
	log.Info("Successfully connected to database.")
	cache := &Cache{
		db:          db,
		initialized: false,
	}
	if err := cache.init(); err != nil {
		log.Panicf("Failed to initialize the database: %e", err)
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

func (c *Cache) WriteEntry(entry *CacheEntry) {
	c.db.Create(entry)
	log.Infof("Wrote cache entry for %s", entry)
}

func (c *Cache) ReadEntry(ordinal int) (*CacheEntry, error) {
	value := new(CacheEntry)
	result := c.db.Where("ordinal = ?", ordinal).First(value)
	if result.Error != nil {
		log.Warningf("Failed to retrieve cache entry: %v", result.Error)
		return nil, result.Error
	}
	log.Debugf("Successfully retrieved cached value for ordinal %d", ordinal)
	return value, nil
}
