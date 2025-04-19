package executor

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/bytecodealliance/wasmtime-go/v31"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
	"github.com/dadamu/contract-wasmvm/internal/contract/runtime"
)

type ContractExecutor struct {
	engine *wasmtime.Engine
}

func NewContractExecutor(
	engine *wasmtime.Engine,
) *ContractExecutor {
	return &ContractExecutor{
		engine: engine,
	}
}

func generateContractId(
	state []byte,
	codeId uint64,
	salt []byte,
) string {
	codeIdBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(codeIdBz, codeId)
	contractId := sha256.Sum256(append(state, append(codeIdBz, salt...)...))
	return base58.Encode(contractId[:])
}

func (ce *ContractExecutor) InitializeContract(
	repository interfaces.IContractRepository,
	state []byte,
	codeId uint64,
	args []byte,
	gasLimit uint64,
) (uint64, string, error) {
	// Get the total contract amount as salt
	amount := repository.GetTotalContractAmount()
	salt := make([]byte, 8)
	binary.LittleEndian.PutUint64(salt, amount)

	contractId := generateContractId(state, codeId, salt)
	err := repository.CreateConctract(codeId, contractId)
	if err != nil {
		return 0, "", err
	}

	remaining, err := ce.RunContract(
		repository,
		state,
		interfaces.ContractMessage{
			Contract: contractId,
			Method:   "init",
			Args:     args,
		},
		gasLimit,
	)

	return remaining, contractId, err
}

func (ce *ContractExecutor) RunContract(
	repository interfaces.IContractRepository,
	state []byte,
	msg interfaces.ContractMessage,
	gasLimit uint64,
) (uint64, error) {
	callbackQueue := callbackqueue.NewCallbackQueue()

	// Enqueue the initial contract call
	// This is the first contract call that will be executed
	callbackQueue.Enqueue(msg)

	for msg, found := callbackQueue.Dequeue(); found; {

		// Run the contract with the current gas limit
		remaining, err := ce.runContract(callbackQueue, repository, state, msg, gasLimit)
		if err != nil {
			return 0, err
		}

		// Update the gas limit for the next contract call
		gasLimit = remaining
	}

	return gasLimit, nil
}

func (ce *ContractExecutor) runContract(
	callbackQueue *callbackqueue.CallbackQueue,
	repository interfaces.IContractRepository,
	state []byte,
	msg interfaces.ContractMessage,
	gasLimit uint64,
) (uint64, error) {
	// Load the contract code from the repository
	module, err := ce.loadContract(repository, msg.Contract)
	if err != nil {
		return 0, err
	}

	if msg.Method == "init" {
		err := repository.TryInitializeContract(msg.Contract)
		if err != nil {
			return 0, fmt.Errorf("failed to initialize contract: %w", err)
		}
	}

	// Execute the contract
	runtime := runtime.NewRuntimeFromModule(ce.engine, callbackQueue, repository, module, state, msg.Contract, gasLimit)
	remaining, err := runtime.Run(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to run contract: %w", err)
	}

	return remaining, nil
}

func (ce *ContractExecutor) loadContract(repository interfaces.IContractRepository, contractId string) (*wasmtime.Module, error) {
	// Load the code from the repository
	code, err := repository.GetContractCodeByContract(contractId)
	if err != nil {
		return nil, err
	}

	// Create a new module from the raw binary
	module, err := wasmtime.NewModule(ce.engine, code)
	if err != nil {
		return nil, err
	}

	return module, nil
}
