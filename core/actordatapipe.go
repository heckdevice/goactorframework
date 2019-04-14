package core

type ActorMessagePipe interface {
	Process(message Message)
	RequestClose()
	Self() ActorBehaviour
	GiveActionableMessage() (ActionableMessage, bool)
	IsAcceptingMessages() bool
}

func (actor *Actor) Process(message Message) {
	actor.dataChan <- message
}
func (actor *Actor) RequestClose() {
	actor.closeChan <- true
}
func (actor *Actor) Self() ActorBehaviour {
	return actor
}

func (actor *Actor) GiveActionableMessage() (ActionableMessage, bool) {
	return actor.internalMessageQueue.Pop()
}

func (actor *Actor) IsAcceptingMessages() bool {
	return actor.isacceptingmessages
}
