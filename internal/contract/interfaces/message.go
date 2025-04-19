package interfaces

type VMMessage interface {
	IsVMMessage()
}

func NewContractMessage(contract, method string, args []byte, sender string) ContractMessage {
	return ContractMessage{
		Contract: contract,
		Method:   method,
		Args:     args,
		Sender:   sender,
	}
}

type ContractMessage struct {
	Contract string
	Method   string
	Args     []byte
	Sender   string
}

func (cm ContractMessage) IsVMMessage() {}

type DeployContractCodeMessage struct {
	Code   []byte
	Sender string
}

func (dccm DeployContractCodeMessage) IsVMMessage() {}

type InitializeContractMessage struct {
	CodeId uint64
	Args   []byte
	Sender string
}

func (icm InitializeContractMessage) IsVMMessage() {}

// ----------------------------------------------------------------------------

type Transaction interface {
	GetGasLimit() uint64
	GetState() []byte
	GetMessages() []VMMessage
}
