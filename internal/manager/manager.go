package manager

import (
	"github.com/bytecodealliance/wasmtime-go/v31"
	"github.com/dadamu/contract-wasmvm/internal/interfaces"
)

type ContractManager struct {
	engine *wasmtime.Engine

	repository interfaces.IContractRepository
}

func NewContractManager(
	engine *wasmtime.Engine,
	repository interfaces.IContractRepository,
) *ContractManager {
	return &ContractManager{
		engine:     engine,
		repository: repository,
	}
}
