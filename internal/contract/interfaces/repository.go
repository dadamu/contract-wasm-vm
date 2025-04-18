package interfaces

type IContractRepository interface {
	// SaveEntity saves the entity with the given key and data in the contract namespace.
	SaveEntity(contractId string, key string, data []byte)

	// LoadEntity loads the entity with the given key in the contract namespace.
	LoadEntity(contractId string, key string) []byte

	// GetContracCode retrieves the code of the contract.
	GetContractCodeByContract(contractId string) ([]byte, error)

	// GetContractCode retrieves the code of the contract by its ID.
	GetContractCodeById(codeId uint64) ([]byte, error)

	// GetContractCodeId retrieves the code ID of the contract.
	CreateConctract(codeId uint64, contractId string) error

	// TryInitializeContract tries to initialize the contract.
	// If the contract is already initialized, it returns an error.
	TryInitializeContract(contractId string) error

	// GetTotalContractAmount retrieves the total amount of the contract.
	GetTotalContractAmount() uint64
}
