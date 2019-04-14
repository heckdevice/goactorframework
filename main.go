package main

import (
	"fmt"
	"github.com/heckdevice/goactorframework/core"
	"github.com/heckdevice/goactorframework/samples"
	"os"
	"os/signal"
	"syscall"
)

var (
	killPill         = make(chan os.Signal)
	terminateProcess = make(chan chan bool)
)

func main() {
	var terminateActorSystemSignal chan bool
	signal.Notify(killPill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGTSTP)
	core.GetDefaultRegistry().InitActorSystem(samples.InitSampleMessageQueue())
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
