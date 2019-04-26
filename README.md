[![CircleCI](https://circleci.com/gh/heckdevice/goactorframework-corelib.svg?style=svg)](https://circleci.com/gh/heckdevice/goactorframework-corelib)
[![Go Report Card](https://goreportcard.com/badge/github.com/heckdevice/goactorframework-corelib)](https://goreportcard.com/report/github.com/heckdevice/goactorframework-corelib)


# goactorframework
A simple actor framework for go

# platfrom version
go version go1.12.1


# Usage

 Get Default Actor system by invoking core.GetDefaultActorSystem()
 
 ActorSystem interface has following features :
 ```
// ActorSystem - Features of actor system
type ActorSystem interface {
	Start(messageQueue chan Message)
	Close(terminateProcess chan bool)
	RegisterActor(actor *Actor, messageType string, handler func(message Message)) error
	UnregisterActor(string) error
	GetActor(actorType string) (ActorMessagePipe, error)
}
 ```
 Start the actor system using Start function which takes the message channel to pick messages from 
 ```
 core.GetDefaultActorSystem().Start(<processing message channle>)
 ```
 
 To write a new actor all we need to do is the following :
 
  - Create a simple actor struct providing ActorType 
  ```
  printActor := core.Actor{ActorType: ActorType}
  ```
  - Using the DefaultActorSystem we register the actor providing a MessageType (string) and its respective handler function
  ```
  err := core.GetDefaultActorSystem().RegisterActor(&printActor, common.ConsolePrint, consolePrint)
  ```
  A handler function can be any function of type 
  ```
  func (message core.Message)
  
  Like in PrintActor
  
  func consolePrint(message core.Message) {
	fmt.Print(fmt.Sprintf("Got Message %v", message))
   }
  ```
  - Spawn the actor in its own routine
  ```
  go printActor.SpawnActor()
  ```
 # References  
  For details refer goactorframework-examples  
  - https://github.com/heckdevice/goactorframework-examples
