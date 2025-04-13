package executor

import (
	"github.com/bytecodealliance/wasmtime-go/v31"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/runtime"
	"github.com/dadamu/contract-wasmvm/internal/interfaces"
)

type ContractExecutor struct {
	engine *wasmtime.Engine

	repository interfaces.IContractRepository
}

func NewContractExecutor(
	engine *wasmtime.Engine,
	repository interfaces.IContractRepository,
) *ContractExecutor {
	return &ContractExecutor{
		engine:     engine,
		repository: repository,
	}
}

func (ce *ContractExecutor) RunContractWithGasLimit(msg interfaces.ContractMessage, gasLimit uint64) error {
	callbackQueue := callbackqueue.NewCallbackQueue()

	// Enqueue the initial contract call
	// This is the first contract call that will be executed
	callbackQueue.Enqueue(msg)

	for msg, found := callbackQueue.Dequeue(); found; {
		remaining, err := ce.runContract(callbackQueue, msg, gasLimit)
		if err != nil {
			return err
		}

		// Update the gas limit for the next contract call
		gasLimit = remaining
	}

	return nil
}

func (ce *ContractExecutor) runContract(callbackQueue *callbackqueue.CallbackQueue, msg interfaces.ContractMessage, gasLimit uint64) (uint64, error) {
	module, err := ce.loadContract(msg.Contract)
	if err != nil {
		return 0, err
	}

	runtime := runtime.NewRuntimeFromModule(callbackQueue, ce.engine, msg.Contract, ce.repository, module, gasLimit)

	// TODO: Add state for run instead of nil
	return runtime.Run(nil, msg)
}

func (ce *ContractExecutor) loadContract(contractId string) (*wasmtime.Module, error) {
	// Load the raw binary module from the repository
	rawModule, err := ce.repository.GetContractRawModule(contractId)
	if err != nil {
		return nil, err
	}

	// Create a new module from the raw binary
	module, err := wasmtime.NewModule(ce.engine, rawModule)
	if err != nil {
		return nil, err
	}

	return module, nil
}
