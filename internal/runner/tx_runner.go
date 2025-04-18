package runner

import (
	"fmt"

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

func (r *TxRunner) RunTransaction(tx interfaces.Transaction) (uint64, error) {
	gasLimit := tx.GetGasLimit()

	msgs := tx.GetMessages()
	state := tx.GetState()

	for _, msg := range msgs {
		gasLimit, err := r.runMessage(state, msg, gasLimit)
		if err != nil {
			return 0, err
		}

		if gasLimit == 0 {
			return 0, fmt.Errorf("gas limit exhausted")
		}
	}

	return gasLimit, nil
}

func (r *TxRunner) runMessage(state []byte, msg interfaces.VMMessage, gasLimit uint64) (uint64, error) {
	switch msg := msg.(type) {

	case interfaces.DeployContractCodeMessage:
		remaining, err := r.runDeployContract(msg, gasLimit)
		if err != nil {
			r.store.Rollback()
			return 0, err
		}
		return remaining, nil

	case interfaces.InitializeContractMessage:
		return 0, fmt.Errorf("initialize contract message not supported")

	case interfaces.ContractMessage:
		remaining, err := r.runContract(state, msg, gasLimit)
		if err != nil {
			r.store.Rollback()
			return 0, err
		}
		return remaining, nil

	default:
		panic("unknown message type")
	}
}

func (r *TxRunner) runDeployContract(msg interfaces.DeployContractCodeMessage, gasLimit uint64) (uint64, error) {
	// Consume gas limit
	consumed := DEPLOY_GAS * uint64(len(msg.Code))
	if consumed > gasLimit {
		return 0, fmt.Errorf("not enough gas limit")
	}

	r.store.StoreContractCode(msg.Code)
	gasLimit -= consumed

	return gasLimit, nil
}

func (r *TxRunner) runContract(state []byte, msg interfaces.ContractMessage, gasLimit uint64) (uint64, error) {
	return r.executor.RunContractWithGasLimit(r.store, state, msg, gasLimit)
}
