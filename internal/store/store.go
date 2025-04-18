package store

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/iavl"
)

type Store struct {
	tree  *iavl.MutableTree
	cache *CacheKVStore
}

func NewStore(tree *iavl.MutableTree) *Store {
	getFn := func(key []byte) []byte {
		value, err := tree.ImmutableTree.Get(key)
		if err != nil {
			panic(err)
		}
		return value
	}

	setFn := func(key, value []byte) {
		_, err := tree.Set(key, value)
		if err != nil {
			panic(err)
		}
	}

	deleteFn := func(key []byte) {
		_, _, err := tree.Remove(key)
		if err != nil {
			panic(err)
		}
	}

	return &Store{
		tree: tree,
		cache: NewCacheKVStore(
			getFn,
			setFn,
			deleteFn,
		),
	}
}

func (s *Store) GetCached() *CacheKVStore {
	return s.cache
}

func (s *Store) SaveVersionWithId(id uint64) ([]byte, error) {
	newVersion := s.tree.WorkingVersion()
	newVersionBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(newVersionBytes, uint64(newVersion))
	s.tree.Set(newVersionKey(id), newVersionBytes)

	hash, _, err := s.tree.SaveVersion()

	return hash, err
}

func (s *Store) GetVersionById(id uint64) (uint64, error) {
	newVersionBytes, err := s.tree.Get(newVersionKey(id))
	if err != nil {
		return 0, err
	}

	if newVersionBytes == nil {
		return 0, fmt.Errorf("version not found for id %d", id)
	}

	return binary.BigEndian.Uint64(newVersionBytes), nil
}
