package interfaces

type ContractMessage struct {
	Contract string
	Method   string
	Args     []byte
	Sender   string
}

func NewContractMessage(contract, method string, args []byte, sender string) ContractMessage {
	return ContractMessage{
		Contract: contract,
		Method:   method,
		Args:     args,
		Sender:   sender,
	}
}

// ----------------------------------------------------------------------------

type Transaction interface {
	GetGasLimit() uint64
	GetState() []byte
	GetContractMessages() []ContractMessage
}
