# goactorframework
A simple actor framework for go

# platfrom version
go version go1.11.2

# how to run
go run main.go

# Usage

 Get Default Registry by invoking core.GetDefaultRegistry()
 
 Registry interface has following features :
 
 type registryInterface interface {
	InitActorSystem(messageQueue chan Message)
	Close()
	RegisterActor(actor *Actor, messageType string, handler func(message Message)) error
	UnregisterActor(string) error
	GetActor(actorType string) (ActorMessagePipe, error)
 }
 
 Initialize the actor system using InitActorSystem function which takes the incoming message channel
 
 In sample examples this is provided by the InitSampleMessageQueue function
 
 Please review main.go code to see the usage
