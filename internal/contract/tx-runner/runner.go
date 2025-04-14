package contract

import (
	"github.com/dadamu/contract-wasmvm/internal/contract/executor"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
)

type TxRunner struct {
	executor executor.ContractExecutor
}

func NewTxRunner(
	executor executor.ContractExecutor,
) *TxRunner {
	return &TxRunner{
		executor: executor,
	}
}

func (r *TxRunner) RunTransaction(tx interfaces.Transaction) error {
	gasLimit := tx.GetGasLimit()
	for _, msg := range tx.GetContractMessages() {
		remaining, err := r.executor.RunContractWithGasLimit(msg, gasLimit)
		if err != nil {
			return err
		}

		// Update the gas limit for the next contract call
		gasLimit = remaining
	}

	return nil
}
