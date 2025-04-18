package store

import (
	"fmt"
	"strconv"

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

// --------------------------------------------------------------

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

func (ck *CacheKVStore) GetContractCodeByContract(
	contractId string,
) ([]byte, error) {
	key := newContractModuleKey(contractId)
	if !ck.has(key) {
		return nil, fmt.Errorf("contract does not exist: %s", contractId)
	}

	codeKey := ck.get(key)
	return ck.get(codeKey), nil
}

func (ck *CacheKVStore) GetContractCodeById(
	codeId uint64,
) ([]byte, error) {
	key := newContractCodeKey(codeId)
	code := ck.get(key)
	if code == nil {
		return nil, fmt.Errorf("contract code not found for id: %d", codeId)
	}
	return code, nil
}

func (ck *CacheKVStore) CreateConctract(
	codeId uint64,
	contractId string,
) error {
	key := newContractModuleKey(contractId)
	if ck.has(key) {
		return fmt.Errorf("contract already exists: %s", contractId)
	}
	ck.set(key, newContractCodeKey(codeId))
	return nil
}

func (ck *CacheKVStore) TryInitializeContract(
	contractId string,
) error {
	key := newContractInitializedKey(contractId)
	if ck.has(key) {
		return fmt.Errorf("contract already initialized: %s", contractId)
	}
	ck.set(key, []byte{0x1})
	return nil
}

func (ck *CacheKVStore) GetTotalContractAmount() uint64 {
	key := []byte(CONTRACT_NEXT_CODE_ID_KEY)
	if !ck.has(key) {
		return 0
	}
	nextCodeIdBtyes := ck.get(key)
	totalAmount, err := strconv.ParseUint(string(nextCodeIdBtyes), 10, 64)
	if err != nil {
		panic("failed to parse next code id")
	}

	return totalAmount
}

func (ck *CacheKVStore) StoreContractCode(
	code []byte,
) {
	nextCodeIdBtyes := ck.get([]byte(CONTRACT_NEXT_CODE_ID_KEY))
	if nextCodeIdBtyes == nil {
		nextCodeIdBtyes = []byte("0")
	}

	nextCodeId, err := strconv.ParseUint(string(nextCodeIdBtyes), 10, 64)
	if err != nil {
		panic("failed to parse next code id")
	}

	key := newContractCodeKey(nextCodeId)
	ck.set(key, code)
	ck.set([]byte(CONTRACT_NEXT_CODE_ID_KEY), []byte(strconv.FormatUint(nextCodeId+1, 10)))
}

// --------------------------------------------------------------

func (ck *CacheKVStore) has(key []byte) bool {
	if cached, ok := ck.cache[string(key)]; ok {
		return !cached.deleted
	}

	return ck.getFn(key) != nil
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
