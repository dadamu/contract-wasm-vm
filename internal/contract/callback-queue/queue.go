package callbackqueue

import (
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
)

type CallbackQueue struct {
	container []interfaces.ContractMessage
	head      int
}

func NewCallbackQueue() *CallbackQueue {
	return &CallbackQueue{
		container: make([]interfaces.ContractMessage, 0),
	}
}
func (q *CallbackQueue) Enqueue(msg interfaces.ContractMessage) {
	q.container = append(q.container, msg)
}

func (q *CallbackQueue) Dequeue() (interfaces.ContractMessage, bool) {
	if q.IsEmpty() {
		return interfaces.ContractMessage{}, false
	}

	msg := q.container[q.head]
	q.head += 1
	return msg, true
}

func (q *CallbackQueue) IsEmpty() bool {
	return q.head >= len(q.container)
}

func (q *CallbackQueue) All() []interfaces.ContractMessage {
	return q.container
}
