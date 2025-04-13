package callbackqueue

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
	container []CallbackMessage
}

func NewCallbackQueue() *CallbackQueue {
	return &CallbackQueue{
		container: make([]CallbackMessage, 0),
	}
}
func (q *CallbackQueue) Enqueue(msg CallbackMessage) {
	q.container = append(q.container, msg)
}

func (q *CallbackQueue) Dequeue() (CallbackMessage, bool) {
	if len(q.container) == 0 {
		return CallbackMessage{}, false
	}
	msg := q.container[0]
	q.container = q.container[1:]
	return msg, true
}

func (q *CallbackQueue) IsEmpty() bool {
	return len(q.container) == 0
}
