package core

// ActorMessagePipe - The actors' data processing interface
type ActorMessagePipe interface {
	Process(message Message)
	AckClose()
	RequestClose()
	Self() ActorBehaviour
	GiveActionableMessage() (ActionableMessage, bool)
	IsAcceptingMessages() bool
}

// Process - This puts the messages to be processed into the actors data channel
func (actor *Actor) Process(message Message) {
	actor.dataChan <- message
}

// AckClose - Acknowledgment by actor for a close request
func (actor *Actor) AckClose() {
	actor.closeChan <- true
}

// RequestClose - Sends a request to close the actor to actors' data channel
func (actor *Actor) RequestClose() {
	actor.Self().getDataChan() <- Message{MessageType: KILLPILL}
}

// Self - Returns ActorBehaviour interface instance of the actor
func (actor *Actor) Self() ActorBehaviour {
	return actor
}

// GiveActionableMessage - Returns the next actionable messages from actors' message queue
func (actor *Actor) GiveActionableMessage() (ActionableMessage, bool) {
	return actor.internalMessageQueue.Pop()
}

// IsAcceptingMessages - Checks if the actor is accepting message for processing
func (actor *Actor) IsAcceptingMessages() bool {
	return actor.isAcceptingMessages
}
