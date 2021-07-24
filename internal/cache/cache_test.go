package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache("postgres://fibo:averysecurepasswordshouldgohere@localhost:15432/fibo")
	defer func() {
		assert.NoError(t, cache.Close())
	}()
}

func TestReadWriteEntry(t *testing.T) {
	cache := NewCache("postgres://fibo:averysecurepasswordshouldgohere@localhost:15432/fibo")
	defer func() {
		assert.NoError(t, cache.Close())
	}()
	// Test writing some entries
	assert.NoError(t, cache.Write(0, 0))
	assert.NoError(t, cache.Write(1, 1))
	assert.NoError(t, cache.Write(2, 1))
	// Test reading the values back
	v, err := cache.Read(0)
	assert.Equal(t, uint64(0), v)
	assert.NoError(t, err)
	v, err = cache.Read(1)
	assert.Equal(t, uint64(1), v)
	assert.NoError(t, err)
	v, err = cache.Read(2)
	assert.Equal(t, uint64(1), v)
	assert.NoError(t, err)
}
