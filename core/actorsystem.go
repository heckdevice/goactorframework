package core

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
)

var (
	actorSys actorSystem
	mutex    sync.Mutex
)

func init() {
	actorSys = actorSystem{Name: "DefaultActorSystem"}
	actorSys.registeredActorsPipe = make(map[string]ActorMessagePipe)
	actorSys.ActorCloseAcked = make(chan bool)
}

type actorSystem struct {
	registeredActorsPipe map[string]ActorMessagePipe
	Name                 string
	ActorCloseAcked      chan bool
}

func GetDefaultRegistry() registryInterface {
	return &actorSys
}

type registryInterface interface {
	InitActorSystem(messageQueue chan Message)
	Close(ackClose chan bool)
	RegisterActor(actor *Actor, messageType string, handler func(message Message)) error
	UnregisterActor(string) error
	GetActor(actorType string) (ActorMessagePipe, error)
}

// RegisterActor - Registers a bare-bone actor to the actor system
// Minimum requirement for an actor to qualify for registration is to have
// its type defined and have at-least one message handler
func (actorSys *actorSystem) RegisterActor(actor *Actor, messageType string, handler func(message Message)) error {
	if actor == nil || len(strings.TrimSpace(actor.ActorType)) == 0 {
		return errors.New(fmt.Sprintf("Invalid actor %v", actor))
	}
	if actorFound, OK := actorSys.registeredActorsPipe[actor.ActorType]; OK {
		return errors.New(fmt.Sprintf("Actor %v is already registerd", actorFound.Self().Type()))
	}
	mutex.Lock()
	actor.id = actor.ActorType + "-" + uuid.New().String()
	actor.handlers = make(map[string]func(Message))
	actor.handlers[messageType] = handler
	actor.internalMessageQueue = make(messageStack, 0, 0)
	actor.dataChan = make(chan Message, 10)
	actor.closeChan = make(chan bool)
	actorSys.registeredActorsPipe[actor.Type()] = actor
	actor.isacceptingmessages = true
	actor.owner = actorSys
	mutex.Unlock()
	return nil
}

func (actorSys *actorSystem) UnregisterActor(actorType string) error {
	if len(strings.TrimSpace(actorType)) == 0 {
		return errors.New("actorType can not be empty")
	}
	if actorFound, OK := actorSys.registeredActorsPipe[actorType]; OK {
		actorFound.RequestClose()
	} else {
		return errors.New(fmt.Sprintf("Actor %v is not registerd", actorType))
	}
	return nil
}

func (actorSys *actorSystem) GetActor(actorType string) (ActorMessagePipe, error) {
	if actorFound, OK := actorSys.registeredActorsPipe[actorType]; OK {
		return actorFound, nil

	}
	return nil, errors.New(fmt.Sprintf("Actor %v is not registerd", actorType))
}

func validateMessage(message Message) error {
	//TODO implement basic validation for message mode and nil checks
	return nil
}

func (actorSys *actorSystem) Close(ackClose chan bool) {
	for _, actor := range actorSys.registeredActorsPipe {
		actor.RequestClose()
	}
	go func(*actorSystem, chan bool) {
		noOfActorsAcked := 0
		log.Printf("Waiting for %v actors to ack close", len(actorSys.registeredActorsPipe))
		for {
			select {
			case <-actorSys.ActorCloseAcked:
				noOfActorsAcked = noOfActorsAcked + 1
				log.Printf("Number of actors to ack close is %v", noOfActorsAcked)
				if noOfActorsAcked == len(actorSys.registeredActorsPipe) {
					log.Println("All actors acked close request")
					return
				}
			}
		}
		ackClose <- true
	}(actorSys, ackClose)
}

func (actorSys *actorSystem) AckActorClosed() {
	actorSys.ActorCloseAcked <- true
}

func (actorSys *actorSystem) InitActorSystem(messageQueue chan Message) {
	go actorSys.actOnMessages()
	go actorSys.startDispatcher(messageQueue)
}
func (actorSys *actorSystem) startDispatcher(incomingMessages chan Message) {
	for {
		select {
		case message := <-incomingMessages:
			err := validateMessage(message)
			if err != nil {
				log.Printf("Invalid message by actor %v, rejecting it, please re-post a valid message", message.Sender.ActorType)
			} else {
				switch message.Mode {
				case Unicast:
					sendToActor, error := actorSys.GetActor(message.UnicastTo.ActorType)
					if error != nil {
						log.Printf("Actor %v not found to process message %v", message.UnicastTo.ActorType, message)
					} else {
						if sendToActor.IsAcceptingMessages() {
							sendToActor.Process(message)
						} else {
							log.Printf("!!!Actor %v is no longer accepting messages, please re-post for processing later!!!", message.Sender.ActorType)
						}
					}
					//TODO Broadcast Support pending
				case Broadcast:
					log.Println("!!!Broadcast feature not yet implemented for all messages!!!")
				}
			}
		}
	}
}

func (actorSys *actorSystem) actOnMessages() {
	for {
		for actorType, actor := range actorSys.registeredActorsPipe {
			if actionableMessage, OK := actor.GiveActionableMessage(); OK {
				log.Printf("Processing message for actor %v", actorType)
				actionableMessage.Handler(actionableMessage.Message)
			}
		}
	}
}
