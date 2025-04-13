package callbackqueue

import "container/list"

type CallbackMessage struct {
	Contract string
	Method   string
	Args     []byte
}

func NewCallbackMessage(contract, method string, args []byte) CallbackMessage {
	return CallbackMessage{
		Contract: contract,
		Method:   method,
		Args:     args,
	}
}

// ---------------------------------------------------------

type CallbackQueue struct {
	container *list.List
}

func NewCallbackQueue() *CallbackQueue {
	return &CallbackQueue{
		container: list.New(),
	}
}
func (q *CallbackQueue) Enqueue(msg CallbackMessage) {
	q.container.PushBack(msg)
}

func (q *CallbackQueue) Dequeue() (CallbackMessage, bool) {
	if q.IsEmpty() {
		return CallbackMessage{}, false
	}

	front := q.container.Front()
	q.container.Remove(front)
	return front.Value.(CallbackMessage), true
}

func (q *CallbackQueue) IsEmpty() bool {
	return q.container.Len() == 0
}
