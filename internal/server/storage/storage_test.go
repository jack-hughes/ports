package storage

import (
	"fmt"
	types "github.com/jack-hughes/ports/internal"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"

	"testing"
)

var someErr = fmt.Errorf("some-error")

func TestStorage_NewStorage(t *testing.T) {
	t.Run("test new successfully", func(t *testing.T) {
		store := NewStorage(zap.NewNop())
		assert.Equal(t, Store{
			log: zap.NewNop(),
			db:  types.InMemStore{Ports: make(map[string]types.Port)},
		}, store)
	})
}

func TestStorage_Update(t *testing.T) {
	t.Run("test update", func(t *testing.T) {
		store := NewStorage(zap.NewNop())
		result := store.Update(getTestInternalPort("aaa"))
		assert.IsType(t, types.Port{}, result)
		assert.Equal(t, getTestInternalPort("aaa"), result)
	})
}

func TestStorage_Get(t *testing.T) {
	t.Run("test get non-existent port", func(t *testing.T) {
		id := "test-port-id"
		store := NewStorage(zap.NewNop())
		result, err := store.Get(id)
		require.Error(t, err)
		assert.Equal(t, fmt.Errorf("could not find port with id: %v", id), err)
		assert.Equal(t, types.Port{}, result)
	})
	t.Run("test get existing port", func(t *testing.T) {
		id := "aaa"
		store := NewStorage(zap.NewNop())
		store.Update(getTestInternalPort(id))
		result, err := store.Get(id)
		assert.NoError(t, err)
		assert.Equal(t, getTestInternalPort(id), result)
	})
}

func TestStorage_List(t *testing.T) {
	t.Run("no records should be an empty array", func(t *testing.T) {
		store := NewStorage(zap.NewNop())
		var list []types.Port
		result := store.List()
		assert.Equal(t, list, result)
	})
	t.Run("test get existing port", func(t *testing.T) {
		store := NewStorage(zap.NewNop())
		store.Update(getTestInternalPort("aaa"))
		store.Update(getTestInternalPort("bbb"))

		result := store.List()
		assert.Equal(t, 2, len(result))
		assert.Equal(t, getTestInternalPort("aaa"), result[0])
		assert.Equal(t, getTestInternalPort("bbb"), result[1])
	})
}

func getTestInternalPort(id string) types.Port {
	return types.Port{
		ID:          id,
		Name:        "bbb",
		City:        "ccc",
		Country:     "ddd",
		Alias:       []string{"1", "2", "3"},
		Regions:     []string{"4", "5", "6"},
		Coordinates: []float32{1.2, 1.3, 1.4},
		Province:    "eee",
		Timezone:    "fff",
		Unlocs:      []string{"7", "8", "9"},
		Code:        "ggg",
	}
}
