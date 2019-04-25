package core

import (
	"errors"
	"fmt"
	"log"
	"time"
)

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
func (actor *Actor) RegisterMessageHandler(messageType string, handler func(message Message)) error {
	if _, OK := actor.handlers[messageType]; OK {
		return errors.New(fmt.Sprintf("Handler for message type %v is already registered for actor %v", messageType, actor.ActorType))
	} else {
		mutex.Lock()
		actor.handlers[messageType] = handler
		mutex.Unlock()
	}
	return nil
}

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
func (actor *Actor) Type() string {
	return actor.ActorType
}

//*************************** Instance methods ***************************
func (actor *Actor) HasMessages() bool {
	return actor.internalMessageQueue.Len() != 0
}

func (actor *Actor) ScheduleActionableMessage(am *ActionableMessage) {
	actor.internalMessageQueue.Push(*am)
}
func (actor *Actor) StopAcceptingMessages() {
	actor.isAcceptingMessages = false
}

func (actor *Actor) NoOfMessagesInQueue() int {
	return actor.internalMessageQueue.Len()
}

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
