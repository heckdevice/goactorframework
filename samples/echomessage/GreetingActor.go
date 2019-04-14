package echomessage

import (
	"fmt"
	"log"
	"tracpilf/core"
	"tracpilf/samples"
)

const (
	ActorType      = "GreetingActor"
	MessageTypeHI  = "HI"
	MessageTypeBYE = "BYE"
)

func InitActor() {
	greetingActor := core.Actor{ActorType: ActorType}
	err := core.GetDefaultActorSystem().RegisterActor(&greetingActor, MessageTypeHI, greetHI)
	if err != nil {
		log.Panic(fmt.Sprintf("Error while registering actor %v. Details : %v", greetingActor.ActorType, err.Error()))
	}
	greetingActor.RegisterMessageHandler(MessageTypeBYE, greetBye)
	greetingActor.RegisterMessageHandler(samples.ConsolePrint, consolePrint)
	go greetingActor.SpawnActor()
}

func greetHI(message core.Message) {
	fmt.Print(fmt.Sprintf("Hi there %v, i got %v", message.Sender.ActorType, message.Payload))
}

func greetBye(message core.Message) {
	fmt.Print(fmt.Sprintf("Adios %v !!!", message.Sender.ActorType))
}

func consolePrint(message core.Message) {
	fmt.Print(fmt.Sprintf("Echo : %v", message))
}
