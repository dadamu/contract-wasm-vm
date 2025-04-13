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

func (ce *ContractExecutor) RunContractWithGasLimit(contractId string, method string, args []byte, gasLimit uint64) error {
	callbackQueue := callbackqueue.NewCallbackQueue()

	// Enqueue the initial contract call
	// This is the first contract call that will be executed
	callbackQueue.Enqueue(callbackqueue.NewCallbackMessage(contractId, method, args))

	for callback, found := callbackQueue.Dequeue(); found; {
		remaining, err := ce.runContract(callbackQueue, callback.Contract, callback.Method, callback.Args, gasLimit)
		if err != nil {
			return err
		}

		// Update the gas limit for the next contract call
		gasLimit = remaining
	}

	return nil
}

func (ce *ContractExecutor) runContract(callbackQueue *callbackqueue.CallbackQueue, contractId string, method string, args []byte, remainingGas uint64) (uint64, error) {
	module, err := ce.loadContract(contractId)
	if err != nil {
		return 0, err
	}

	runtime := runtime.NewRuntimeFromModule(callbackQueue, ce.engine, contractId, ce.repository, module, remainingGas)

	// TODO: Add state for run instead of nil
	return runtime.Run(method, nil, args)
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
