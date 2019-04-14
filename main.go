package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tracpilf/core"
	"tracpilf/samples"
	"tracpilf/samples/echomessage"
	"tracpilf/samples/printmessage"
)

var (
	killPill         = make(chan os.Signal)
	terminateProcess = make(chan chan bool)
	messageQueue     = make(chan core.Message, 10)
)

func PumpMessages() {
	dummySender := core.ActorReference{"ActorSystem"}
	for {
		messageQueue <- core.Message{MessageType: samples.ConsolePrint,
			Mode:        core.Broadcast,
			Sender:      &dummySender,
			Payload:     map[string]interface{}{"data": rand.Int()},
			BroadcastTo: []*core.ActorReference{&core.ActorReference{ActorType: printmessage.ActorType}, &core.ActorReference{ActorType: echomessage.ActorType}}}
		messageQueue <- core.Message{MessageType: samples.ConsolePrint, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		messageQueue <- core.Message{MessageType: samples.ConsolePrint, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: printmessage.ActorType}}
		messageQueue <- core.Message{MessageType: echomessage.MessageTypeHI, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		messageQueue <- core.Message{MessageType: echomessage.MessageTypeBYE, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		time.Sleep(time.Second * 3)
	}
}
func main() {
	var terminateActorSystemSignal chan bool
	signal.Notify(killPill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGTSTP)
	printmessage.InitActor()
	echomessage.InitActor()
	go PumpMessages()
	core.GetDefaultActorSystem().InitActorSystem(messageQueue)
	for {
		select {
		case <-killPill:
			fmt.Println(fmt.Sprintf("\n\n****** - Shutting down due to SIGTERM - ******"))
			//core.ProcessKillPill(terminateActorSystemSignal)
			fmt.Println("\n\n****** - SIGTERM Processed - ******")
			return
		case <-terminateActorSystemSignal:
			//Define function in Actor for graceful shutdown
		case <-terminateProcess:
			fmt.Println("\n\n****** - Terminate process command received - ******")
			//core.ProcessKillPill(terminateActorSystemSignal)
			return
		}
	}

}
