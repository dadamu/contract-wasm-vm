package runner

import (
	"fmt"

	checker "github.com/dadamu/contract-wasmvm/internal/code-checker"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/executor"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
	"github.com/dadamu/contract-wasmvm/internal/store"
)

const DEPLOY_GAS = uint64(0)

type TxRunner struct {
	executor executor.ContractExecutor
	store    *store.CacheKVStore
}

func NewTxRunner(
	executor executor.ContractExecutor,
	store *store.CacheKVStore,
) *TxRunner {
	return &TxRunner{
		executor: executor,
		store:    store,
	}
}

func (r *TxRunner) RunTransaction(tx interfaces.Transaction) (uint64, []*callbackqueue.CallbackQueue, error) {
	gasLimit := tx.GetGasLimit()

	msgs := tx.GetMessages()
	state := tx.GetState()

	queues := make([]*callbackqueue.CallbackQueue, 0)
	for _, msg := range msgs {
		remaining, callbackQueue, err := r.runMessage(state, msg, gasLimit)
		if err != nil {
			return 0, nil, err
		}

		gasLimit = remaining
		queues = append(queues, callbackQueue)
	}

	return gasLimit, queues, nil
}

func (r *TxRunner) runMessage(state []byte, msg interfaces.VMMessage, gasLimit uint64) (uint64, *callbackqueue.CallbackQueue, error) {
	switch msg := msg.(type) {

	case interfaces.DeployContractCodeMessage:
		remaining, err := r.deployContract(msg, gasLimit)
		if err != nil {
			r.store.Rollback()
			return 0, nil, err
		}
		return remaining, nil, nil

	case interfaces.InitializeContractMessage:
		remaining, callbackQueue, err := r.executor.CreateContract(
			r.store,
			state,
			msg.CodeId,
			msg.Args,
			gasLimit,
		)
		if err != nil {
			r.store.Rollback()
			return 0, callbackQueue, err
		}
		return remaining, callbackQueue, nil

	case interfaces.ContractMessage:
		remaining, callbackQueue, err := r.executor.RunContract(r.store, state, msg, gasLimit)
		if err != nil {
			r.store.Rollback()
			return 0, callbackQueue, err
		}
		return remaining, callbackQueue, nil

	default:
		panic("unknown message type")
	}
}

func (r *TxRunner) deployContract(msg interfaces.DeployContractCodeMessage, gasLimit uint64) (uint64, error) {
	// Consume gas limit
	consumed := DEPLOY_GAS * uint64(len(msg.Code))
	if consumed > gasLimit {
		return 0, fmt.Errorf("not enough gas limit")
	}

	isUndeterminstic, err := checker.ContainUndeterminsticOps(msg.Code)
	if err != nil {
		return 0, fmt.Errorf("failed to check code: %w", err)
	}

	if isUndeterminstic {
		return 0, fmt.Errorf("code contains unsupported operations")
	}

	r.store.StoreContractCode(msg.Code)
	gasLimit -= consumed

	return gasLimit, nil
}
