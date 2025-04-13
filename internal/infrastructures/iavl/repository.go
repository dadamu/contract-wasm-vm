package iavl

import (
	"github.com/cosmos/iavl"
	"github.com/dadamu/contract-wasmvm/internal/interfaces"
)

const CONTRACT_ENTITY_PREFIX = "contracts/entities"
const CONTRACT_MODULE_PREFIX = "contracts/modules"

func newContractEntityKey(contractId string, id string) []byte {
	return []byte(CONTRACT_ENTITY_PREFIX + "/" + contractId + "/" + id)
}

func newContractModuleKey(contractId string) []byte {
	return []byte(CONTRACT_MODULE_PREFIX + "/" + contractId)
}

// -------------------------------------------------------------------------

var _ interfaces.IContractRepository = (*IAVLRepository)(nil)

type IAVLRepository struct {
	tree *iavl.MutableTree
}

func NewIAVLRepository(contractId string, tree *iavl.MutableTree) *IAVLRepository {
	return &IAVLRepository{
		tree: tree,
	}
}

func (r *IAVLRepository) SaveEntity(id string, contractId string, data []byte) error {
	key := newContractEntityKey(contractId, id)
	_, err := r.tree.Set(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *IAVLRepository) LoadEntity(contractId string, id string) ([]byte, error) {
	key := newContractEntityKey(contractId, id)
	value, err := r.tree.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (r *IAVLRepository) GetContractRawModule(contractId string) ([]byte, error) {
	key := newContractModuleKey(contractId)
	value, err := r.tree.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}
