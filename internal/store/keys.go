package store

import (
	"fmt"
)

const CONTRACT_ENTITY_PREFIX = "contracts/entities"
const CONTRACT_MODULE_PREFIX = "contracts/modules"
const CONTRACT_INITIALIZED_PREFIX = "contracts/initialized"

const CONTRACT_CODE_PREFIX = "contracts/codes"
const CONTRACT_NEXT_CODE_ID_KEY = "contracts/next_code_id"

const VERSION_MAP_PREFIX = "version"

func newContractEntityKey(contractId string, id string) []byte {
	key := fmt.Sprintf("%s/%s/%s", CONTRACT_INITIALIZED_PREFIX, contractId, id)
	return []byte(key)
}

func newContractInitializedKey(contractId string) []byte {
	key := fmt.Sprintf("%s/%s", CONTRACT_ENTITY_PREFIX, contractId)
	return []byte(key)
}

func newContractModuleKey(contractId string) []byte {
	key := fmt.Sprintf("%s/%s", CONTRACT_MODULE_PREFIX, contractId)
	return []byte(key)
}

func newContractCodeKey(codeId uint64) []byte {
	key := fmt.Sprintf("%s/%d", CONTRACT_CODE_PREFIX, codeId)
	return []byte(key)
}

func parseContractCodeKey(key []byte) (uint64, error) {
	var codeId uint64
	_, err := fmt.Sscanf(string(key), "%s/%d", &codeId)
	if err != nil {
		return 0, err
	}
	return codeId, nil
}

func newVersionKey(id uint64) []byte {
	key := fmt.Sprintf("%s/%d", VERSION_MAP_PREFIX, id)
	return []byte(key)
}
