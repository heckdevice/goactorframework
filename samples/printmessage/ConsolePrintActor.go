package printmessage

import (
	"fmt"
	"log"
	"tracpilf/core"
	"tracpilf/samples"
)

const (
	ActorType = "PrintActor"
)

func InitActor() {
	printActor := core.Actor{ActorType: ActorType}
	err := core.GetDefaultActorSystem().RegisterActor(&printActor, samples.ConsolePrint, consolePrint)
	if err != nil {
		log.Panic(fmt.Sprintf("Error while registering actor %v. Details : %v", printActor.ActorType, err.Error()))
	}
	go printActor.SpawnActor()
}

func consolePrint(message core.Message) {
	fmt.Print(fmt.Sprintf("Got Message %v", message))
}
