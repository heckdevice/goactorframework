[![CircleCI](https://circleci.com/gh/heckdevice/goactorframework.svg?style=svg)](https://circleci.com/gh/heckdevice/goactorframework)
[![Go Report Card](https://goreportcard.com/badge/github.com/heckdevice/goactorframework)](https://goreportcard.com/report/github.com/heckdevice/goactorframework)


# goactorframework
A simple actor framework for go

# platfrom version
go version go1.12.1

# how to run
go run main.go

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
 
 In sample examples this is provided by the InitSampleMessageQueue function
 
 ```
 func InitSampleMessageQueue() chan core.Message {
	printmessage.InitActor()
	echomessage.InitActor()
	go pumpMessages()
	return messageQueue
}
 ```
 There we initialze two sample actor of type "PrintActor" and "GreetingActor"
 
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
 Please review main.go and any of the actors in "samples" directory to see the detailed design and usage
