package cache

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/stretchr/testify/assert"
	pg "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var database string = "fibo_test"
var connString string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + database})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	connString = fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s", resource.GetPort("5432/tcp"), database)
	if err := pool.Retry(func() error {
		_, err := gorm.Open(pg.Open(connString), &gorm.Config{
			// This turns off the default logging which is too verbose for records that don't exist
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run the test
	retCode := m.Run()

	// Purge the container
	// Can't be deferred because os.Exit doesn't respect defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(retCode)
}

func TestCreateCache(t *testing.T) {
	cache := NewCache(connString)
	defer func() {
		assert.NoError(t, cache.Close())
	}()
}

func TestReadWriteEntry(t *testing.T) {
	cache := NewCache(connString)
	defer func() {
		assert.NoError(t, cache.Close())
	}()
	// Test writing some entries
	assert.NoError(t, cache.Write(0, fibonacci.NewNumber(0)))
	assert.NoError(t, cache.Write(1, fibonacci.NewNumber(1)))
	assert.NoError(t, cache.Write(2, fibonacci.NewNumber(1)))
	// Test reading the values back
	v, err := cache.Read(0)
	assert.Equal(t, fibonacci.NewNumber(0), v)
	assert.NoError(t, err)
	v, err = cache.Read(1)
	assert.Equal(t, fibonacci.NewNumber(1), v)
	assert.NoError(t, err)
	v, err = cache.Read(2)
	assert.Equal(t, fibonacci.NewNumber(1), v)
	assert.NoError(t, err)
}
