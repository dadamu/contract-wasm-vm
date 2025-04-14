package store

import (
	"fmt"
)

const CONTRACT_ENTITY_PREFIX = "contracts/entities"
const CONTRACT_MODULE_PREFIX = "contracts/modules"
const VERSION_MAP_PREFIX = "version"

func newContractEntityKey(contractId string, id string) []byte {
	key := fmt.Sprintf("%s/%s/%s", CONTRACT_ENTITY_PREFIX, contractId, id)
	return []byte(key)
}

func newContractModuleKey(contractId string) []byte {
	key := fmt.Sprintf("%s/%s", CONTRACT_MODULE_PREFIX, contractId)
	return []byte(key)
}

func newVersionKey(id uint64) []byte {
	key := fmt.Sprintf("%s/%d", VERSION_MAP_PREFIX, id)
	return []byte(key)
}
