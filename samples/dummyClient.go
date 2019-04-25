package samples

import (
	"math/rand"
	"time"

	"github.com/heckdevice/goactorframework/core"
	"github.com/heckdevice/goactorframework/samples/common"
	"github.com/heckdevice/goactorframework/samples/echomessage"
	"github.com/heckdevice/goactorframework/samples/printmessage"
)

var (
	messageQueue = make(chan core.Message, 10)
)

func InitSampleMessageQueue() chan core.Message {
	printmessage.InitActor()
	echomessage.InitActor()
	go pumpMessages()
	return messageQueue
}

func pumpMessages() {
	dummySender := core.ActorReference{}
	dummySender.ActorType = "ActorSystem"
	for {
		messageQueue <- core.Message{MessageType: common.ConsolePrint,
			Mode:        core.Broadcast,
			Sender:      &dummySender,
			Payload:     map[string]interface{}{"data": rand.Int()},
			BroadcastTo: []*core.ActorReference{{ActorType: printmessage.ActorType}, {ActorType: echomessage.ActorType}}}
		messageQueue <- core.Message{MessageType: common.ConsolePrint, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		messageQueue <- core.Message{MessageType: common.ConsolePrint, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: printmessage.ActorType}}
		messageQueue <- core.Message{MessageType: echomessage.MessageTypeHI, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		messageQueue <- core.Message{MessageType: echomessage.MessageTypeBYE, Mode: core.Unicast, Sender: &dummySender, Payload: map[string]interface{}{"data": rand.Int()}, UnicastTo: &core.ActorReference{ActorType: echomessage.ActorType}}
		time.Sleep(time.Millisecond * 50)
	}
}
