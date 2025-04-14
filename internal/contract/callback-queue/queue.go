package callbackqueue

import (
	"container/list"

	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
)

type CallbackQueue struct {
	container *list.List
}

func NewCallbackQueue() *CallbackQueue {
	return &CallbackQueue{
		container: list.New(),
	}
}
func (q *CallbackQueue) Enqueue(msg interfaces.ContractMessage) {
	q.container.PushBack(msg)
}

func (q *CallbackQueue) Dequeue() (interfaces.ContractMessage, bool) {
	if q.IsEmpty() {
		return interfaces.ContractMessage{}, false
	}

	front := q.container.Front()
	q.container.Remove(front)
	return front.Value.(interfaces.ContractMessage), true
}

func (q *CallbackQueue) IsEmpty() bool {
	return q.container.Len() == 0
}
