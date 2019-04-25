package core

// ActionableMessage - Coalesce message with its registered handler
type ActionableMessage struct {
	Message
	Handler func(message Message)
}

type messageStack []*ActionableMessage

func (b *messageStack) Clear() {
	if len(*b) == 0 {
		return
	}

	for i := range *b {
		(*b)[i] = nil
	}
	*b = (*b)[:0]
}

func (b *messageStack) Peek() (v ActionableMessage, ok bool) {
	l := b.Len()
	if l > 0 {
		ok = true
		v = *(*b)[l-1]
	}
	return
}

func (b *messageStack) Push(v ActionableMessage) {
	mutex.Lock()
	*b = append(*b, &v)
	mutex.Unlock()
}

func (b *messageStack) Pop() (v ActionableMessage, ok bool) {
	mutex.Lock()
	l := b.Len()
	if l > 0 {
		l--
		ok = true
		v = *(*b)[l]
		(*b)[l] = nil
		*b = (*b)[:l]
	}
	mutex.Unlock()
	return
}

func (b *messageStack) Len() int {
	return len(*b)
}
