package iavl

import (
	"github.com/cosmos/iavl"
)

const CONTRACT_PREFIX = "contracts"

type IAVLRepository struct {
	tree *iavl.MutableTree
}

func NewIAVLRepository(contractId string, tree *iavl.MutableTree) *IAVLRepository {
	return &IAVLRepository{
		tree: tree,
	}
}

func (r *IAVLRepository) SaveEntity(id string, contractId string, data []byte) error {
	key := []byte(CONTRACT_PREFIX + "/" + contractId + "/" + id)
	_, err := r.tree.Set(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *IAVLRepository) LoadEntity(contractId string, id string) ([]byte, error) {
	key := []byte(CONTRACT_PREFIX + "/" + contractId + "/" + id)
	value, err := r.tree.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}
