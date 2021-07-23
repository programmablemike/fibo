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
	assert.NoError(t, cache.Init())
	assert.Equal(t, true, cache.Initialized())
}

func TestReadWriteEntry(t *testing.T) {
	cache := NewCache("fibo", "averysecurepasswordshouldgohere", "localhost:15432", "fibo")
	defer func() {
		assert.NoError(t, cache.Close())
	}()
	assert.NoError(t, cache.Init())
	// Test writing some entries
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 0, Result: 1}))
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 1, Result: 1}))
	assert.NoError(t, cache.WriteEntry(&CacheEntry{Ordinal: 2, Result: 2}))
	// Test reading the values back
	res, err := cache.ReadEntry(0)
	assert.EqualValues(t, &CacheEntry{Ordinal: 0, Result: 1}, res)
	assert.NoError(t, err)
	res, err = cache.ReadEntry(1)
	assert.EqualValues(t, &CacheEntry{Ordinal: 1, Result: 1}, res)
	assert.NoError(t, err)
	res, err = cache.ReadEntry(2)
	assert.EqualValues(t, &CacheEntry{Ordinal: 2, Result: 2}, res)
	assert.NoError(t, err)
}
