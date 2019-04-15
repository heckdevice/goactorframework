package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/heckdevice/goactorframework/core"
	"github.com/heckdevice/goactorframework/samples"
)

var (
	killPill         = make(chan os.Signal)
	terminateProcess = make(chan bool)
)

func main() {
	signal.Notify(killPill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGTSTP)
	oncomingMessages := samples.InitSampleMessageQueue()
	core.GetDefaultRegistry().InitActorSystem(oncomingMessages)
	for {
		select {
		case <-killPill:
			fmt.Println(fmt.Sprintf("\n\n******--- Shutting down due to SIGTERM ---******"))
			core.GetDefaultRegistry().Close(terminateProcess)
		case <-terminateProcess:
			fmt.Println(fmt.Sprintf("\n\n******--- Actor system is stopped, exiting ---******"))
			return
		}
	}

}
