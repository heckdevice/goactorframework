# goactorframework
A simple actor framework for go

# platfrom version
go version go1.11.2

# how to run
go run main.go

# Usage

 Get Default Registry by invoking core.GetDefaultRegistry()
 
 Registry interface has following features :
 ```
 type registryInterface interface {
	InitActorSystem(messageQueue chan Message)
	Close()
	RegisterActor(actor *Actor, messageType string, handler func(message Message)) error
	UnregisterActor(string) error
	GetActor(actorType string) (ActorMessagePipe, error)
 }
 ```
 Initialize the actor system using InitActorSystem function which takes the incoming message channel
 
 as in main.go
 ```
 core.GetDefaultRegistry().InitActorSystem(samples.InitSampleMessageQueue())
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
  - Using the DefaultRegistry register the actor providing a MessageType (string) and its respective handler function
  ```
  err := core.GetDefaultRegistry().RegisterActor(&printActor, common.ConsolePrint, consolePrint)
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
