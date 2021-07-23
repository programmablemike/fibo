package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache("fibo", "averysecurepasswordshouldgohere", "localhost:15432", "fibo")
	defer func() {
		assert.NoError(t, cache.Close())
	}()
}

func TestReadWriteEntry(t *testing.T) {
	cache := NewCache("fibo", "averysecurepasswordshouldgohere", "localhost:15432", "fibo")
	defer func() {
		assert.NoError(t, cache.Close())
	}()
	// Test writing some entries
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 0, Result: 0}))
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 1, Result: 1}))
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 2, Result: 1}))
	// Test reading the values back
	res, err := cache.ReadEntry(0)
	assert.Equal(t, 0, res.Ordinal)
	assert.Equal(t, 0, res.Result)
	assert.NoError(t, err)
	res, err = cache.ReadEntry(1)
	assert.Equal(t, 1, res.Ordinal)
	assert.Equal(t, 1, res.Result)
	assert.NoError(t, err)
	res, err = cache.ReadEntry(2)
	assert.Equal(t, 2, res.Ordinal)
	assert.Equal(t, 1, res.Result)
	assert.NoError(t, err)
}
