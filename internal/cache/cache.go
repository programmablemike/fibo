package cache

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	pg "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

type CacheEntry struct {
	gorm.Model
	Ordinal uint64 `gorm:"uniqueIndex"` // The fibonacci ordinal N
	Value   uint64 // The fibonacci value
}

func (c CacheEntry) String() string {
	return fmt.Sprintf("CacheEntry<%s %s>", strconv.FormatUint(c.Ordinal, 10), strconv.FormatUint(c.Value, 10))
}

// Cache implements a PostgresDB cache for pre-computed ordinal values
type Cache struct {
	db          *gorm.DB
	initialized bool
}

// NewCache creates a new cache with persistent database connection
func NewCache(dsn string) *Cache {
	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("Failed to connect to database: %e", err)
	}
	log.Info("Successfully connected to database.")
	cache := &Cache{
		db:          db,
		initialized: false,
	}
	if err := cache.init(); err != nil {
		log.Errorf("Failed to initialize the database: %e", err)
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

func (c *Cache) Write(ordinal uint64, value uint64) error {
	entry := &CacheEntry{
		Ordinal: ordinal,
		Value:   value,
	}
	c.db.Create(entry)
	log.Infof("Wrote cache entry for %s", strconv.FormatUint(ordinal, 10))
	return nil
}

func (c *Cache) Read(ordinal uint64) (uint64, error) {
	entry := new(CacheEntry)
	result := c.db.Where("ordinal = ?", ordinal).First(entry)
	if result.Error != nil {
		log.Warningf("Failed to retrieve cache entry: %v", result.Error)
		return 0, result.Error
	}
	log.Debugf("Successfully retrieved cached value for ordinal %s", strconv.FormatUint(ordinal, 10))
	return entry.Value, nil
}
