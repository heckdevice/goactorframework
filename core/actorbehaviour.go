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

func (actor *Actor) SpawnActor() {
	for {
		select {
		case data := <-actor.dataChan:
			switch data.MessageType {
			case KillPill:
				//stop accepting messages
				actor.StopAcceptingMessages()
				//wait till all the messages are processed and the internal message queue is empty
				log.Println(fmt.Sprintf("Actor %v with id %v got KillPill message", actor.ActorType, actor.id))
				for {
					if actor.HasMessages() {
						log.Println(fmt.Sprintf("Actor %v still have messages in pipe, waiting for finishing the process before shutting down", actor.ActorType))
						time.Sleep(time.Millisecond * 500)
					} else {
						close(actor.dataChan)
						actor.RequestClose()
						break
					}
				}
				return
			default:
				//Default behaviour is to delegate the message to the actor pipe for processing
				//as per the registered handlers
				log.Println(fmt.Sprintf("Actor %v with id %v got message", actor.ActorType, actor.id))
				if handlerFound, OK := actor.GetRegisteredHandlers()[data.MessageType]; OK {
					actor.ScheduleActionableMessage(&ActionableMessage{data, handlerFound})
				} else {
					log.Fatalf(fmt.Sprintf("Actor %v has no handler for message type %v, rejecting the message", actor.ActorType, data.MessageType))
				}
			}
		case <-actor.closeChan:
			log.Println(fmt.Sprintf("Actor %v closing down", actor.ActorType))
			close(actor.closeChan)
			return
		}
	}
}
