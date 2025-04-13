package interfaces

type IContractRepository interface {
	// SaveEntity saves the entity with the given key and data in the contract namespace
	SaveEntity(contractId string, key string, data []byte) error
	// LoadEntity loads the entity with the given key in the contract namespace
	LoadEntity(contractId string, key string) ([]byte, error)

	GetContractRawModule(contractId string) ([]byte, error)
}
