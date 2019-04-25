package core

import (
	"fmt"
	"log"
	"time"
)

// ActorBehaviour - Actor features interface
type ActorBehaviour interface {
	RegisterMessageHandler(messageType string, handler func(message Message)) error
	GetRegisteredHandlers() map[string]func(Message)
	getDataChan() chan Message
	setDataChan(dataChan chan Message)
	getCloseChan() chan bool
	setCloseChan(dataChan chan bool)
	Type() string
}

//*************************** ActorBehaviour interface methods ***************************

// RegisterMessageHandler - This enables registering the handler function for a MessageType for an actor
func (actor *Actor) RegisterMessageHandler(messageType string, handler func(message Message)) error {
	if _, OK := actor.handlers[messageType]; OK {
		return fmt.Errorf("handler for message type %v is already registered for actor %v", messageType, actor.ActorType)
	}
	mutex.Lock()
	actor.handlers[messageType] = handler
	mutex.Unlock()
	return nil
}

// GetRegisteredHandlers - Returns a map of all messagetypes and the respective registered handler function
func (actor *Actor) GetRegisteredHandlers() map[string]func(Message) {
	return actor.handlers
}
func (actor *Actor) getDataChan() chan Message {
	return actor.dataChan
}
func (actor *Actor) setDataChan(dataChan chan Message) {
	actor.dataChan = dataChan
}
func (actor *Actor) getCloseChan() chan bool {
	return actor.closeChan
}
func (actor *Actor) setCloseChan(closeChan chan bool) {
	actor.closeChan = closeChan
}

// Type - Returns the Actor type/name
func (actor *Actor) Type() string {
	return actor.ActorType
}

//*************************** Instance methods ***************************

// HasMessages - Returns true if any messages are pending to be processed in the actors' message stack
func (actor *Actor) HasMessages() bool {
	return actor.internalMessageQueue.Len() != 0
}

// ScheduleActionableMessage - This schedules the ActionableMessage for an actor by pushing it into its message stack
func (actor *Actor) ScheduleActionableMessage(am *ActionableMessage) {
	actor.internalMessageQueue.Push(*am)
}

// StopAcceptingMessages - Stops the actor for accepting any messages, this generally needs to be invoked just after de-registering the actor
func (actor *Actor) StopAcceptingMessages() {
	actor.isAcceptingMessages = false
}

// NoOfMessagesInQueue - Returns the number of messages scheduled and pending the the actors' message queue
func (actor *Actor) NoOfMessagesInQueue() int {
	return actor.internalMessageQueue.Len()
}

// SpawnActor - This starts the actors' message processing go routine. For the actor to start accpeting any message and there by processing it this is a mandatory invocation
func (actor *Actor) SpawnActor() {
	for {
		select {
		case data := <-actor.dataChan:
			switch data.MessageType {
			case KILLPILL:
				//stop accepting messages
				log.Println(fmt.Sprintf("Stopping Actor %v with id %v to accept any more messages", actor.ActorType, actor.id))
				actor.StopAcceptingMessages()
				go func(actor *Actor) {
					for {
						if actor.HasMessages() {
							log.Printf("!!!Actor %v still have %v messages in pipe!!!", actor.ActorType, actor.NoOfMessagesInQueue())
							time.Sleep(time.Millisecond * 250)
						} else {
							log.Printf("!!!Actor %v have no more messages in pipe!!!", actor.ActorType)
							actor.AckClose()
							break
						}
					}
				}(actor)
			default:
				//Default behaviour is to delegate the message to the actor pipe for processing
				//as per the registered handlers
				log.Println(fmt.Sprintf("Actor %v with id %v got message", actor.ActorType, actor.id))
				if actor.isAcceptingMessages {
					if handlerFound, OK := actor.GetRegisteredHandlers()[data.MessageType]; OK {
						actor.ScheduleActionableMessage(&ActionableMessage{data, handlerFound})
					} else {
						log.Fatalf(fmt.Sprintf("Actor %v has no handler for message type %v, rejecting the message", actor.ActorType, data.MessageType))
					}
				}
			}
		case <-actor.closeChan:
			log.Println(fmt.Sprintf("Actor %v closing down due to close signal", actor.ActorType))
			actor.owner.AckActorClosed()
			close(actor.dataChan)
			close(actor.closeChan)
			actor.internalMessageQueue.Clear()
			return
		}
	}
}
