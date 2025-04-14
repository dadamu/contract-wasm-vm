package store

import (
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
)

var _ interfaces.IContractRepository = (*CacheKVStore)(nil)

type CacheKVStore struct {
	cache          map[string]cacheValue
	updatedKeyList []string
	updatedKeyMap  map[string]bool

	// Functions to interact with the underlying store
	getFn    func([]byte) []byte
	setFn    func([]byte, []byte)
	deleteFn func([]byte)
}

type cacheValue struct {
	value   []byte
	deleted bool
}

func NewCacheKVStore(
	getFn func([]byte) []byte,
	setFn func([]byte, []byte),
	deleteFn func([]byte),
) *CacheKVStore {
	return &CacheKVStore{
		cache:          make(map[string]cacheValue),
		updatedKeyList: make([]string, 0),
		updatedKeyMap:  make(map[string]bool),

		getFn:    getFn,
		setFn:    setFn,
		deleteFn: deleteFn,
	}
}

func (ck *CacheKVStore) LoadEntity(
	contractId, entityKey string,
) []byte {
	key := newContractEntityKey(contractId, entityKey)

	return ck.get(key)
}

func (ck *CacheKVStore) SaveEntity(
	contractId, entityKey string, data []byte,
) {
	key := newContractEntityKey(contractId, entityKey)
	ck.set(key, data)
}

func (ck *CacheKVStore) GetContractRawModule(
	contractId string,
) []byte {
	key := newContractModuleKey(contractId)
	return ck.get(key)
}

func (ck *CacheKVStore) get(key []byte) []byte {
	if cached, ok := ck.cache[string(key)]; ok {
		if cached.deleted {
			return nil
		}
		return cached.value
	}

	return ck.getFn(key)
}

func (ck *CacheKVStore) set(key, value []byte) {
	keyStr := string(key)
	ck.cache[keyStr] = cacheValue{
		value:   value,
		deleted: false,
	}

	if !ck.updatedKeyMap[keyStr] {
		ck.updatedKeyList = append(ck.updatedKeyList, keyStr)
		ck.updatedKeyMap[keyStr] = true
	}
}

func (ck *CacheKVStore) Rollback() {
	ck.cache = make(map[string]cacheValue)
	ck.updatedKeyList = make([]string, 0)
	ck.updatedKeyMap = make(map[string]bool)
}

func (ck *CacheKVStore) Commit() {
	for _, keyStr := range ck.updatedKeyList {
		cahced := ck.cache[keyStr]
		key := []byte(keyStr)

		if cahced.deleted {
			ck.deleteFn(key)
		} else {
			ck.setFn(key, cahced.value)
		}
	}

	// Clear the cache after committing
	ck.cache = make(map[string]cacheValue)
	ck.updatedKeyList = make([]string, 0)
	ck.updatedKeyMap = make(map[string]bool)
}
