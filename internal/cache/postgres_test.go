package cache

import (
	"log"
	"testing"

	"github.com/ory/dockertest"
	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/stretchr/testify/assert"
)

var database string = "fibo_test"
var connString = "postgres://postgres:secret@localhost:5432/fibo_test"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + database})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
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
