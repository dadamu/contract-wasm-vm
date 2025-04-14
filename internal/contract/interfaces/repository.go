package interfaces

type IContractRepository interface {
	// SaveEntity saves the entity with the given key and data in the contract namespace
	SaveEntity(contractId string, key string, data []byte)

	// LoadEntity loads the entity with the given key in the contract namespace
	LoadEntity(contractId string, key string) []byte

	// GetContractRawModule retrieves the raw binary module of the contract
	GetContractRawModule(contractId string) []byte
}
