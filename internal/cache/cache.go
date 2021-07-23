package cache

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
)

type CacheEntry struct {
	Id      int64 // The primary key
	Ordinal int   // The fibonacci ordinal N
	Result  int   // The fibonacci value
}

func (c CacheEntry) String() string {
	return fmt.Sprintf("CacheEntry<%d %d>", c.Ordinal, c.Result)
}

// Cache implements a PostgresDB cache for pre-computed ordinal values
type Cache struct {
	db          *pg.DB
	initialized bool
}

// NewCache creates a new cache with persistent database connection
func NewCache(user string, password string, addr string, database string) *Cache {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Addr:     addr,
		Database: database,
	})

	return &Cache{
		db:          db,
		initialized: false,
	}
}

// Initialized is for checking if the database schema was initialized
func (c Cache) Initialized() bool {
	return c.initialized
}

// Init the cache database
func (c *Cache) Init() error {
	log.Info("Initializing the database...")
	if c.initialized {
		log.Warning("Cannot re-initiliaze the database... skipping.")
		return nil
	}
	if err := c.initSchema(); err != nil {
		log.Error("Failed to initialize the database schema.")
		return err // TODO: Fail gracefully
	}
	log.Info("Database initialized successfully!")
	c.initialized = true // Mark the database as initialized
	return nil
}

// initSchema creates the table schema
func (c *Cache) initSchema() error {
	models := []interface{}{
		(*CacheEntry)(nil),
	}
	for _, model := range models {
		err := c.db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			log.Errorf("Failed to initialize table for %s", model)
			return err
		}
	}
	log.Info("Successfully initialized the table schemas.")
	return nil
}

func (c *Cache) Close() error {
	log.Info("Closing the database connection.")
	return c.db.Close()
}

func (c *Cache) WriteEntry(entry *CacheEntry) error {
	if _, err := c.db.Model(entry).Insert(); err != nil {
		log.Errorf("Failed to write the cache entry for %s`", entry)
		return err
	}
	log.Infof("Wrote cache entry for %s", entry)
	return nil
}

func (c *Cache) ReadEntry(ordinal int) (*CacheEntry, error) {
	entry := new(CacheEntry)
	if err := c.db.Model(entry).Where("ordinal = ?", ordinal).Select(); err != nil {
		log.Errorf("Failed to retrieve cache entry: %v", err)
		return nil, err
	}
	log.Infof("Successfully retrieved cached value for ordinal %d", ordinal)
	return entry, nil
}
